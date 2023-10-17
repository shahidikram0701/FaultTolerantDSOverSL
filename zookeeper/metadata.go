package zookeeper

import (
	"errors"
	"fmt"
	"sync"
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

func (md *Metadata) DeleteDoneEntries() {
	md.Lock()
	defer md.Unlock()

	for gsn, item := range md.items {
		if item.State == Done {
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
