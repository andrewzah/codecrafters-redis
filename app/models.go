package main

import (
	"sync"
	"time"
)

type AppContext struct {
	StoreMutex    sync.RWMutex
	Store         map[string]RedisValue
	MetadataMutex sync.RWMutex
	Metadata      AppMetadata
}

type AppMetadata struct {
	Role string
}

type RedisValue struct {
	Data string

	// in milliseconds
	Expiry int64
}

func (c *AppContext) InsertData(key, value string, expiryMillis int64) {
	timeMillis := time.Now().UnixMilli()

	Debugf("acquiring lock")
	c.StoreMutex.Lock()
	defer c.StoreMutex.Unlock()

	Debugf("actually inserting data")
	if expiryMillis == -1 {
		c.Store[key] = RedisValue{value, -1}
	} else {
		c.Store[key] = RedisValue{value, timeMillis + expiryMillis}
	}
}

func (c *AppContext) GetData(key string) string {
	timeMillis := time.Now().UnixMilli()

	c.StoreMutex.RLock()
	defer c.StoreMutex.RUnlock()

	val, ok := c.Store[key]
	if !ok {
		return ""
	}

	if val.Expiry != -1 && timeMillis > val.Expiry {
		delete(c.Store, key)
		return ""
	}

	return val.Data
}

type InfoReplication struct {
	role string
}
