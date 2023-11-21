package zookeeper

import (
	context "context"
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/scalog/scalog/logger"

	pb2 "github.com/scalog/scalog/zookeeper/zookeepermetadatapb"
	pb "github.com/scalog/scalog/zookeeper/zookeeperpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Zookeeper struct {
	sync.RWMutex
	consensus *Consensus
	// data structure
	trie *Trie
}

func (zk *Zookeeper) ProbeForMetadataAndUpdate() {
	go zk.consensus.metadata.ProbeAndUpdate()
	go zk.consensus.metadata.TrimMetadata()
}

func (zk *Zookeeper) FetchShardIdFromLog(gsn int64) (int32, error) {
	return zk.consensus.FetchShardId(gsn)
}

func (zk *Zookeeper) GetShardIdForGSNFromMetadata(gsn int64) (int32, error) {
	return zk.consensus.metadata.GetShardId(gsn)
}

func (zk *Zookeeper) GetAllMetadata() ([]*pb2.MetadataListItem, error) {
	return zk.consensus.metadata.GetAllMetadata()
}

func (zk *Zookeeper) GetShardIdMapping(gsn int64) (int32, error) {
	sid, err := zk.consensus.metadata.FetchShardId(gsn)

	if err != nil {
		log.Fatalf("[ Zookeeper ]%v", err)
		return -1, err
	}
	return sid, nil
}

func (zk *Zookeeper) CheckAndUpdateTrie() {
	sleepForIfNoUpdate := time.Duration(int(viper.GetInt("zk-sleep-if-no-update"))) * time.Millisecond
	for {
		zk.Lock()
		nextOpAt := zk.consensus.LSN
		shardId, err := zk.consensus.metadata.FetchShardId(nextOpAt)

		if err != nil {
			log.Printf("[ Zookeeper ][ UpdateTrie ]No new update to apply; %v", err)
			time.Sleep(sleepForIfNoUpdate)
			zk.Unlock()
			continue
		}

		op, err := zk.consensus.ReadFromLog(nextOpAt, shardId)
		if err != nil {
			log.Printf("[ Zookeeper ][ UpdateTrie ]Reading operation from log failed for lsn: %v at sid: %v with error: %v", nextOpAt, shardId, err)
			zk.Unlock()
			continue
		}

		log.Printf("[ Zookeeper ][ UpdateTrie ]Updating trie with operation: %v", op)
		zk.trie.Execute(op)
		zk.consensus.metadata.UpdateEntryState(nextOpAt)
		zk.consensus.IncrementLSN()
		zk.Unlock()
		time.Sleep(sleepForIfNoUpdate)
	}
}

func (zk *Zookeeper) ApplyAllPendingOpsAndReturnRead(latestGSN int64, readOp string) (int, []byte, error) {
	zk.Lock()
	defer zk.Unlock()

	currentLSN := zk.consensus.LSN

	log.Printf("[ Zookeeper ][ ApplyAllPendingOpsAndReturnRead ]Fetching all shards for seq numbers: %v to %v(both included)", currentLSN, latestGSN)

	operations, err := zk.consensus.ReadBulkData(currentLSN, latestGSN)

	if err != nil {
		log.Printf("[ Zookeeper ][ ApplyAllPendingOpsAndReturnRead ]Error reading all the pending ops from seq number: %v to %v", currentLSN, latestGSN)
		return -1, nil, err
	}

	for idx, op := range operations {
		log.Printf("[ Zookeeper ][ ApplyAllPendingOpsAndReturnRead ]Applying operation %v at sequence number: %v", op, currentLSN+int64(idx))
		zk.trie.Execute(op)
	}

	for idx := currentLSN; idx <= latestGSN; idx++ {
		zk.consensus.metadata.UpdateEntryState(int64(idx))
	}

	log.Printf("[ Zookeeper ][ ApplyAllPendingOpsAndReturnRead ]Applied all pending ops")
	zk.consensus.UpdateLSN(latestGSN + 1)

	// log.Printf("[ Zookeeper ][ ApplyAllPendingOpsAndReturnRead ]Trie: \n%v", zk.trie.PrintTrie(zk.trie.Root, 0))

	ver, data, err := zk.trie.ExecuteGet(readOp)

	log.Printf("[ Zookeeper ][ ApplyAllPendingOpsAndReturnRead ]got read for operation %v: %v", readOp, data)

	if err != nil {
		log.Printf("[ Zookeeper ][ ApplyAllPendingOpsAndReturnRead ]Error: %v", err)
		return -1, nil, err
	}

	return ver, []byte(data), nil

}

type ZKServer struct {
	pb.UnimplementedZooKeeperServer
}

var (
	zkState *Zookeeper
)

func ZKInit() {
	consensusModule, err := NewConsensusModule()

	if err != nil {
		log.Fatalf("[ ZooKeeper ]Error initialising zookeeper: %v", err)
		return
	}

	zkState = &Zookeeper{
		consensus: consensusModule,
		trie:      &Trie{Root: &TrieNode{Value: "/"}},
	}
	zkPort := int32(viper.GetInt("zk-port"))
	zid := int32(viper.GetInt("zid"))
	asyncTrieUpdate := viper.GetBool("async-trie-update")

	zkMetadataPort := int32(viper.GetInt("zk-metadata-port"))

	go StartZKMetadataServer(zkMetadataPort + zid)
	zkState.ProbeForMetadataAndUpdate()
	if asyncTrieUpdate {
		go zkState.CheckAndUpdateTrie()
	}
	StartZKServer(zkPort + zid)
}

func StartZKServer(port int32) {
	log.Printf("[ Zookeeper ]Starting Zookeeper server on port %v", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterZooKeeperServer(s, &ZKServer{})

	log.Printf("[ Zookeeper ]Zookeeper Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[ Zookeeper ]Failed to serve: %v", err)
	}
}

func (zk *ZKServer) CreateZNode(ctx context.Context, znode *pb.ZNode) (*pb.Path, error) {
	// write
	op := "CREATE"
	path := znode.GetPath()
	data := string(znode.GetData()[:])

	record := fmt.Sprintf("%v::%v::%v", op, path, data)
	gsn, shard, err := zkState.consensus.WriteToLog(record)

	if err != nil {
		return nil, err
	}

	log.Printf("[ Zookeeper ][ CreateZNode ]Write successful. Wrote data into scalog at gsn: %v shard: %v", gsn, shard)
	return &pb.Path{Path: path}, nil
	// return &pb.Path{Path: fmt.Sprintf("%v/%v", gsn, shard)}, nil // This is just for testing. Return the path instead, ie, uncomment the above line.
}

func (zk *ZKServer) DeleteZNode(ctx context.Context, path *pb.Path) (*pb.Empty, error) {
	return nil, nil
}

func (zk *ZKServer) ExistsZNode(ctx context.Context, path *pb.Path) (*pb.Stat, error) {
	return nil, nil
}

func (zk *ZKServer) GetZNode(ctx context.Context, path *pb.Path) (*pb.ZNode, error) {
	p := path.GetPath()
	log.Printf("[ Zookeeper ][ GetZNode ]Path: %v", p)

	op := "GET"
	record := fmt.Sprintf("%v::%v", op, p)
	gsn, shard, err := zkState.consensus.WriteToLog(record)

	if err != nil {
		log.Printf("[ Zookeeper ][ GetZNode ]Error writing the read operation %v; err: %v", op, err)
		gsn, shard, err = zkState.consensus.WriteToLog(record)
		if err != nil {
			log.Printf("[ Zookeeper ][ GetZNode ]Error writing the read operation %v; err: %v", op, err)
			return nil, err
		}
	}

	log.Printf("[ Zookeeper ][ GetZNode ][ GetZNode ]Inserted Read operation at gsn: %v and shard: %v", gsn, shard)

	ver, data, err := zkState.ApplyAllPendingOpsAndReturnRead(gsn, record)

	if err != nil {
		log.Printf("[ Zookeeper ][ GetZNode ]Error reading data: %v", err)
		return nil, err
	}

	return &pb.ZNode{
		Path: p,
		Data: data,
		Stat: &pb.Stat{
			Version: int32(ver),
		},
	}, nil

}

func (zk *ZKServer) SetZNode(ctx context.Context, request *pb.SetZNodeRequest) (*pb.Stat, error) {
	return nil, nil
}

func GetZNodeChildren(ctx context.Context, path *pb.Path) (*pb.GetZNodeChildrenResponse, error) {
	return nil, nil
}
