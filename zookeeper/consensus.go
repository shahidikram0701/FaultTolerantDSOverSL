package zookeeper

import (
	"errors"
	"fmt"
	"log"

	client "github.com/scalog/scalog/client"
	"github.com/scalog/scalog/pkg/address"
	"github.com/spf13/viper"
)

type Consensus struct {
	LSN int32 // last sequence number in the log that is applied

	scalogClient *client.Client
	metadata     *Metadata
}

func initScalogClient() (*client.Client, error) {
	numReplica := int32(viper.GetInt("data-replication-factor"))
	discPort := uint16(viper.GetInt("disc-port"))
	discAddr := address.NewLocalDiscAddr(discPort)
	dataPort := uint16(viper.GetInt("data-port"))
	dataAddr := address.NewLocalDataAddr(numReplica, dataPort)
	zkClient, err := client.NewClient(dataAddr, discAddr, numReplica)
	if err != nil {
		return nil, err
	}
	return zkClient, nil
}

func initMetadata() *Metadata {
	metadata := &Metadata{
		items: make(map[int64]MetadataListItem),
	}

	return metadata
}

func NewConsensusModule() (*Consensus, error) {
	scalogClient, err := initScalogClient()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[ ConsensusModule ]Error initializing scalogClient: %v\n", err))
	}

	metadata := initMetadata()

	consensusModule := &Consensus{
		LSN:          0,
		scalogClient: scalogClient,
		metadata:     metadata,
	}
	log.Printf("[ ConsensusModule ]Initialised consensus module successfully")
	return consensusModule, nil
}

func (c *Consensus) WriteToLog(record string) (int64, int32, error) {
	if record == "" {
		return -1, -1, errors.New("nothing to write")
	}

	gsn, shard, err := c.scalogClient.AppendOne(record)

	if err != nil {
		return -1, -1, errors.New(fmt.Sprintf("Error writing to the log: %v\n", err))
	}
	// write the gsn, shard map to the metadata
	if err := c.metadata.Insert(gsn, shard); err != nil {
		log.Printf("[ ConsensusModule ][ WriteToLog ]Insert %v->%v to the metadata failed: %v", gsn, shard, err)
	}

	return gsn, shard, nil
}

func (c *Consensus) ReadFromLog(gsn int64, sid int32) (string, error) {
	if gsn < 0 || sid < 0 {
		return "", errors.New("Invalid gsn or sid")
	}
	rid := int32(0) // TODO: support multiple replica read
	record, err := c.scalogClient.Read(gsn, sid, rid)

	if err != nil {
		return "", errors.New(fmt.Sprintf("Error reading from reading log: %v\n", err))
	}

	return record, nil
}
