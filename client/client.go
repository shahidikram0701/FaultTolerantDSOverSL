package client

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/scalog/scalog/data/datapb"
	disc "github.com/scalog/scalog/discovery"
	"github.com/scalog/scalog/discovery/discpb"
	log "github.com/scalog/scalog/logger"
	"google.golang.org/grpc"
)

// ShardingPolicy determines which records are appended to which shards.
type ShardingPolicy func(view *disc.View, record string) (int32, int32)

type Client struct {
	clientID       int32
	numReplica     int32
	nextCSN        int32
	nextGSN        int32
	viewID         int32
	view           *disc.View
	viewC          chan *discpb.View
	appendC        chan *datapb.Record
	ackC           chan *datapb.Ack
	subC           chan *datapb.Record
	shardingPolicy ShardingPolicy

	discAddr   string
	discConn   *grpc.ClientConn
	discClient *discpb.Discovery_DiscoverClient
	discMu     sync.Mutex

	dataConn map[int32]*grpc.ClientConn
	dataMu   sync.Mutex

	localRun bool
}

func NewLocalClient(discAddr string, numReplica int32) (*Client, error) {
	c, err := NewClient(discAddr, numReplica)
	if err != nil {
		return nil, err
	}
	c.localRun = true
	return c, nil
}

func NewClient(discAddr string, numReplica int32) (*Client, error) {
	c := &Client{
		clientID:   generateClientID(),
		numReplica: numReplica,
		nextCSN:    -1,
		nextGSN:    0,
		viewID:     0,
		localRun:   false,
	}
	c.shardingPolicy = NewDefaultShardingPolicy(numReplica).Shard
	c.viewC = make(chan *discpb.View)
	c.appendC = make(chan *datapb.Record)
	c.ackC = make(chan *datapb.Ack)
	c.subC = make(chan *datapb.Record)
	err := c.UpdateDiscovery(discAddr)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func generateClientID() int32 {
	seed := rand.NewSource(time.Now().UnixNano())
	return rand.New(seed).Int31()
}

func (c *Client) UpdateDiscovery(addr string) error {
	c.discMu.Lock()
	defer c.discMu.Unlock()
	if c.discConn != nil {
		c.discConn.Close()
		c.discConn = nil
	}
	c.discAddr = addr
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return fmt.Errorf("Dial discovery at %v failed: %v", addr, err)
	}
	c.discConn = conn
	discClient := discpb.NewDiscoveryClient(conn)
	callOpts := []grpc.CallOption{}
	discDiscoveryClient, err := discClient.Discover(context.Background(), &discpb.Empty{}, callOpts...)
	if err != nil {
		return fmt.Errorf("Create replicate client to %v failed: %v", addr, err)
	}
	c.discClient = &discDiscoveryClient
	return nil
}

func (c *Client) connDataServer(shard, replica int32) (*grpc.ClientConn, error) {
	globalReplicaID := shard*c.numReplica + replica
	var addr string
	if c.localRun {
		addr = fmt.Sprintf("127.0.0.1:%v", 23282+globalReplicaID)
	} else {
		// TODO request kubedns for replica address
	}
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	if conn, ok := c.dataConn[globalReplicaID]; ok && conn != nil {
		c.dataConn[globalReplicaID].Close()
		delete(c.dataConn, globalReplicaID)
	}
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("Dial dataovery at %v failed: %v", addr, err)
	}
	c.dataConn[globalReplicaID] = conn
	return conn, nil
}

func (c *Client) getDataServerConn(shard, replica int32) (*grpc.ClientConn, error) {
	globalReplicaID := shard*c.numReplica + replica
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	if conn, ok := c.dataConn[globalReplicaID]; ok && conn != nil {
		return conn, nil
	}
	return c.connDataServer(shard, replica)
}

func (c *Client) disconnDataServer(shard, replica int32) error {
	globalReplicaID := shard*c.numReplica + replica
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	if conn, ok := c.dataConn[globalReplicaID]; ok && conn != nil {
		c.dataConn[globalReplicaID].Close()
		delete(c.dataConn, globalReplicaID)
	}
	return nil
}

func (c *Client) Start() {
	go c.processView()
	go c.processAppend()
	go c.processAck()
}

func (c *Client) processView() {
	for v := range c.viewC {
		c.view.Update(v)
	}
}

func (c *Client) processAppend() {
}

func (c *Client) doProcessAppend(appendC chan *datapb.Record, shard, replica int32) {
	conn, err := c.getDataServerConn(shard, replica)
	if err != nil {
		// TODO handle the errors
		log.Errorf("%v", err)
		return
	}
	dataClient := datapb.NewDataClient(conn)
	dataAppendClient, err := dataClient.Append(context.Background())
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	for r := range appendC {
		err := dataAppendClient.Send(r)
		if err != nil {
			log.Errorf("%v", err)
			return
		}
	}
}

func (c *Client) processAck() {
	for r := range c.ackC {
		_ = r
	}
}

func (c *Client) getNextClientSN() int32 {
	csn := atomic.AddInt32(&c.nextCSN, 1)
	return csn
}

func (c *Client) Append(record string) (int32, int32, error) {
	r := &datapb.Record{
		ClientID: c.clientID,
		ClientSN: c.getNextClientSN(),
		Record:   record,
	}
	c.appendC <- r
	return 0, 0, nil
}

func (c *Client) SetShardingPolicy(p ShardingPolicy) {
	c.shardingPolicy = p
}