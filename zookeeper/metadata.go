package zookeeper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb2 "github.com/scalog/scalog/zookeeper/zookeepermetadatapb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Metadata struct {
	sync.RWMutex
	items map[int64]MetadataListItem // gsn -> shard-id
}

type EntryStatus int

const (
	Pending EntryStatus = iota
	Done
)

type MetadataListItem struct {
	shardId int32
	State   EntryStatus
	GSN     int64
}

func (metadataListItem MetadataListItem) String() string {
	return fmt.Sprintf("[shard:%v || gsn:%v || State:%v\n]", metadataListItem.shardId, metadataListItem.GSN, metadataListItem.State)
}

func (md *Metadata) String() string {
	md.RLock()
	defer md.RUnlock()
	res := "\n"

	for _, item := range md.items {
		res += fmt.Sprintf("%v\n", item)
	}

	return res
}

type ZKMetadataServer struct {
	pb2.UnimplementedZooKeeperServer
}

func (md *Metadata) DeleteDoneEntries() {
	md.Lock()
	defer md.Unlock()

	trimThreshold := int(viper.GetInt("metadata-trim-threshold"))

	// only when the size of the metadata table is such that it has atlast
	// {trimThreshold} of elements, we do the trim. Calculate the bytes it
	// occupies. Can be upto a few hundred kilobytes. We need to keep completed
	// entries for a while until other nodes get the mapping too. If deleted too
	// early then the nodes will have to hit the log more often which can be expensive

	if len(md.items) < trimThreshold {
		log.Printf("[ Zookeeper Metadata ][ TrimMetadata ] No Trim: Not reached Trim Threshold of %v. Current metadata size: %v", trimThreshold, len(md.items))
		return
	}

	for gsn, item := range md.items {
		if item.State == Done {
			log.Printf("[ Zookeeper Metadata ][ DeleteDoneEntries ] Removing %v", item)
			delete(md.items, gsn)
		}
	}
}

func (md *Metadata) GetShardId(gsn int64) (int32, error) {
	md.RLock()
	defer md.RUnlock()

	if item, found := md.items[gsn]; found {
		return item.shardId, nil
	}

	return -1, errors.New(fmt.Sprintf("[ Metadata ]shard-id not found for gsn: %d", gsn))
}

func (md *Metadata) GetAllMetadata() ([]*pb2.MetadataListItem, error) {
	md.RLock()
	defer md.RUnlock()

	// Initialize a slice to store the metadata items
	metadataList := make([]*pb2.MetadataListItem, 0, len(md.items))

	for _, item := range md.items {
		metadataList = append(metadataList, &pb2.MetadataListItem{
			ShardId: item.shardId,
			GSN:     item.GSN,
		})
	}

	return metadataList, nil
}

func (md *Metadata) Insert(gsn int64, sid int32) error {
	md.Lock()
	defer md.Unlock()

	if _, ok := md.items[gsn]; ok {
		return errors.New(fmt.Sprintf("[ Metadata ][ Insert ]Duplicate gsn: %v", gsn))
	}

	md.items[gsn] = MetadataListItem{
		shardId: sid,
		GSN:     gsn,
		State:   Pending,
	}

	return nil
}

func (md *Metadata) ProbeAndUpdate() {
	interval := time.Duration(int(viper.GetInt("metadata-probe-frequency"))) * time.Millisecond
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	zkNodesIp := viper.GetStringSlice("zk-servers")
	zids_ := viper.Get("zids").([]interface{})
	var zids []int

	for _, zid_ := range zids_ {
		zids = append(zids, zid_.(int))
	}

	myIp := GetOutboundIP()
	i := 0

	for {
		select {
		case <-ticker.C:
			zkMetadataPort := int32(viper.GetInt("zk-metadata-port"))
			zkIpAddr := zkNodesIp[i]
			zid := zids[i]

			i = ((i + 1) % len(zkNodesIp))

			if zkIpAddr == myIp {
				zkIpAddr = zkNodesIp[i]
				zid = zids[i]
				i = ((i + 1) % len(zkNodesIp))
			}

			serverAddress := fmt.Sprintf("%v:%v", zkIpAddr, zkMetadataPort+int32(zid))
			log.Printf("[ ZookKeeper Metadata ] Requestng for metadata from %v", serverAddress)

			conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("[ ZookKeeper Metadata ] Failed to connect to ZkMetadata server: %v", err)
			}
			client := pb2.NewZooKeeperClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			response, err := client.FetchAll(ctx, &pb2.FetchAllRequest{})

			if err == nil {
				log.Printf("[ ZooKeeper Metadata ] Recieved List! Updating metadata")
				md.Update(response.GetMetadataList())
			}
			conn.Close()
			cancel()

		}
	}
}

func (md *Metadata) Update(recievedMetadataList []*pb2.MetadataListItem) {
	log.Printf("[ ZooKeeper Metadata ][ Update ]Current metadata: %v", md.items)
	// log.Printf("[ ZooKeeper Metadata ][ Update ]Recieved metadata: %v", recievedMetadataList)
	for _, item := range recievedMetadataList {
		err := md.Insert(item.GetGSN(), item.GetShardId())

		if err == nil {
			log.Printf("[ ZookKeeper Metadata ][ Update ] Inserted: %v -> %v", item.GetGSN(), item.GetShardId())
		} else {
			log.Printf("[ ZooKeeper Metadata ][ Update ] Already had mapping for %v -> %v", item.GetGSN(), item.GetShardId())
		}
	}
}

// Hadnt found shardid mapping for this gsn during batch updates
// Now probe everyone and find shard of a specific gsn
func (md *Metadata) FetchShardIdFromPeers(gsn int64) (int32, error) {
	// List of Zookeeper server IP addresses
	zkNodesIp := viper.GetStringSlice("zk-servers")
	zids_ := viper.Get("zids").([]interface{})
	var zids []int

	for _, zid_ := range zids_ {
		zids = append(zids, zid_.(int))
	}
	myIp := GetOutboundIP()

	resultCh := make(chan int32)

	var wg sync.WaitGroup

	for idx, ip := range zkNodesIp {
		if ip == myIp {
			continue
		}
		wg.Add(1)
		go func(zkNodeIP string, zid int32) {
			defer wg.Done()
			zkMetadataPort := int32(viper.GetInt("zk-metadata-port"))
			serverAddress := fmt.Sprintf("%v:%v", zkNodeIP, zkMetadataPort+zid)
			conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("[ ZookKeeper Metadata ][ FetchShardIdFromPeers ] Failed to connect to ZkMetadata server: %v", err)
			}
			client := pb2.NewZooKeeperClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			defer conn.Close()

			response, err := client.FetchOne(ctx, &pb2.FetchOneRequest{Gsn: gsn})

			if err == nil && response != nil {
				// Successfully fetched the shard ID, send it to the result channel
				resultCh <- response.GetShardId()
			}
		}(ip, int32(zids[idx]))
	}

	// goroutine to wait for all requests to finish
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Return the first result received from any Zookeeper server
	for shardID := range resultCh {

		return shardID, nil
	}

	return -1, errors.New(fmt.Sprintf("[ ZookKeeper Metadata ][ FetchShardIdFromPeers ] Didnt find shard for gsn %v from any peer", gsn))
}

func (md *Metadata) FetchShardIdFromLog(gsn int64) (int32, error) {
	sid, err := zkState.FetchShardIdFromLog(gsn)

	if err != nil {
		log.Printf("[ ZookKeeper Metadata ][ FetchShardIdFromLog ] Didnt find the gsn: %v in the log", gsn)
		return -1, errors.New(fmt.Sprintf("[ ZookKeeper Metadata ][ FetchShardIdFromLog ] Didnt find the gsn: %v in the log", gsn))
	}

	return sid, err
}

// If the local metadata store doesnt have the shard mapping
// then probe the peers to find the specific gsn
// if found from the peers then update local store and return
// else probe the shared log and find the shard mapping
func (md *Metadata) FetchShardId(gsn int64) (int32, error) {
	sid, err := md.GetShardId(gsn)
	if err == nil {
		return sid, nil
	}

	sid, err = md.FetchShardIdFromPeers(gsn)
	if err == nil {
		md.Insert(gsn, sid)
		return sid, nil
	}

	sid, err = md.FetchShardIdFromLog(gsn)

	if err == nil {
		md.Insert(gsn, sid)
		return sid, nil
	}

	return -1, errors.New("Invalid gsn")
}

func (md *Metadata) TrimMetadata() {
	interval := time.Duration(int(viper.GetInt("metadata-trim-frequency"))) * time.Millisecond
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("[ Zookeeper Metadata ][ TrimMetadata ] Trying a metadata trim")
			md.DeleteDoneEntries()
		}
	}
}

func (md *Metadata) UpdateEntryState(gsn int64) {
	md.Lock()
	defer md.Unlock()

	if item, ok := md.items[gsn]; ok {
		item.State = Done

		md.items[gsn] = item
	}
}

func StartZKMetadataServer(port int32) {
	log.Printf("[ Zookeeper Metadata ]Starting Zookeeper Metadata server on port %v", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb2.RegisterZooKeeperServer(s, &ZKMetadataServer{})

	log.Printf("[ Zookeeper Metadata ]Zookeeper Metadata Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[ Zookeeper Metadata ]Failed to serve: %v", err)
	}
}

func (zkMetadata *ZKMetadataServer) FetchOne(ctx context.Context, request *pb2.FetchOneRequest) (*pb2.FetchOneResponse, error) {
	reqGSN := request.GetGsn()

	shardId, err := zkState.GetShardIdForGSNFromMetadata(reqGSN)

	if err != nil {
		log.Printf("[ Zookeeper Metadata ]GSN not present here: %v", err)

		return nil, err
	}

	return &pb2.FetchOneResponse{
		ShardId: shardId,
	}, nil
}

func (zkMetadata *ZKMetadataServer) FetchAll(ctx context.Context, request *pb2.FetchAllRequest) (*pb2.FetchAllResponse, error) {

	metadataItems, err := zkState.GetAllMetadata()

	if err != nil {
		log.Printf("[ Zookeeper Metadata ]Error fetching all metadata: %v", err)

		return nil, err
	}

	return &pb2.FetchAllResponse{
		MetadataList: metadataItems,
	}, nil
}
