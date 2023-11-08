package zookeeper

import (
	context "context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	pb2 "github.com/scalog/scalog/zookeeper/zookeepermetadatapb"
	pb "github.com/scalog/scalog/zookeeper/zookeeperpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Zookeeper struct {
	sync.RWMutex
	consensus *Consensus
	// data strucure
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
	}
	zkPort := int32(viper.GetInt("zk-port"))
	zid := int32(viper.GetInt("zid"))

	zkMetadataPort := int32(viper.GetInt("zk-metadata-port"))

	go StartZKMetadataServer(zkMetadataPort + zid)
	zkState.ProbeForMetadataAndUpdate()
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
	op := "write"
	path := znode.GetPath()
	data := string(znode.GetData()[:])

	record := fmt.Sprintf("%v::%v::%v", op, path, data)
	gsn, shard, err := zkState.consensus.WriteToLog(record)

	if err != nil {
		return nil, err
	}

	// return &pb.Path{Path: path}, nil
	log.Printf("[ Zookeeper ][ CreateZNode ]Write successful")
	return &pb.Path{Path: fmt.Sprintf("%v/%v", gsn, shard)}, nil // This is just for testing. Return the path instead, ie, uncomment the above line.
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
	// below is only for testing
	brokenPath := strings.Split(p, "/")
	gsn, _ := strconv.Atoi(brokenPath[0])
	sid, _ := strconv.Atoi(brokenPath[1])

	log.Printf("[ Zookeeper ][ GetZNode ]GSN=%v; SID=%v\n", gsn, sid)

	record, err := zkState.consensus.ReadFromLog(int64(gsn), int32(sid))

	if err != nil {
		return nil, err
	}
	return &pb.ZNode{
		Path: p,
		Data: []byte(record),
	}, nil
}

func (zk *ZKServer) SetZNode(ctx context.Context, request *pb.SetZNodeRequest) (*pb.Stat, error) {
	return nil, nil
}

func GetZNodeChildren(ctx context.Context, path *pb.Path) (*pb.GetZNodeChildrenResponse, error) {
	return nil, nil
}
