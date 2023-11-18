package zookeeper

import (
	"errors"
	"fmt"

	log "github.com/scalog/scalog/logger"

	client "github.com/scalog/scalog/client"
	"github.com/scalog/scalog/pkg/address"
	"github.com/spf13/viper"
)

type Consensus struct {
	LSN int64 // last sequence number in the log that is applied

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

func (c *Consensus) IncrementLSN() {
	c.LSN++
}

func (c *Consensus) UpdateLSN(newLSN int64) {
	c.LSN = newLSN
}

func (c *Consensus) WriteToLog(record string) (int64, int32, error) {
	if record == "" {
		return -1, -1, errors.New("nothing to write")
	}

	log.Printf("[ ConsensusModule ]Appending %v to the log", record)
	gsn, shard, err := c.scalogClient.AppendOne(record)
	log.Printf("[ ConsensusModule ]Appended %v to the log; (%v, %v)", record, gsn, shard)

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

func (c *Consensus) FetchShardId(gsn int64) (int32, error) {
	if gsn < 0 {
		return -1, errors.New("Invalid GSN")
	}
	rid := int32(0)
	numShards := int(viper.GetInt("shards"))

	for i := 0; i < numShards; i++ {
		sid := int32(i)
		record, err := c.scalogClient.Read(gsn, sid, rid)

		if record != "" && err == nil {
			return sid, nil
		}
	}

	return -1, errors.New(fmt.Sprintf("[ ConsensusModule ][ FetchShardId ] Failed to find shard if for the gsn: %v", gsn))
}

func (c *Consensus) ReadBulkData(gsn1 int64, gsn2 int64) ([]string, error) {
	// Get Shard Ids of all the sequence numbers from gsn1 to gsn2 both included
	shardIds := make(map[int64]int32)

	// For all the sequence numbers from gsn1 to gsn2 call FetchShardIds and populate the shardIds array
	for gsn := gsn1; gsn <= gsn2; gsn++ {
		shardID, err := c.metadata.FetchShardId(gsn)
		if err != nil {
			log.Printf("[ Consensus ][ ReadBulkData ] ShardId fetch for gsn: %v Failed; %v", gsn, err)
		}
		shardIds[gsn] = shardID
	}

	// Initialize a slice to accumulate the records
	records := make([]string, 0)

	// For each of the gsn between gsn1 and gsn2, use the ReadFromLog method to fetch the record and append it to the records slice
	for gsn := gsn1; gsn <= gsn2; gsn++ {
		shardID, exists := shardIds[gsn]
		if !exists || shardID == -1 {
			// return nil, errors.New(fmt.Sprintf("[ Consensus ][ ReadBulkData ]Shard ID not found for GSN: %v", gsn))
			log.Printf("[ Consensus ][ ReadBulkData ]Shard ID not found for GSN: %v", gsn)
			continue
		}

		record, err := c.ReadFromLog(gsn, shardID)
		if err != nil {
			log.Fatalf("[ Consensus ][ ReadBulkData ]Error reading record for gsn: %v, sid: %v; %v", gsn, shardID, err)
			return nil, err
		}

		// Append the record to the records slice
		records = append(records, record)
	}

	return records, nil
}
