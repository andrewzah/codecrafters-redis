package main

import (
	"sync"
	"time"
)

type MemoryStore struct {
	HashMap   map[string]RedisValue
	DataMutex sync.RWMutex
}

type RedisValue struct {
	Data string

	// in milliseconds
	Expiry int64
}

func (ms MemoryStore) InsertData(key, value string, expiryMillis int64) {
	timeMillis := time.Now().UnixMilli()

	Debugf("acquiring lock")
	ms.DataMutex.Lock()
	defer ms.DataMutex.Unlock()

	Debugf("actually inserting data")
	if expiryMillis == -1 {
		ms.HashMap[key] = RedisValue{value, -1}
	} else {
		ms.HashMap[key] = RedisValue{value, timeMillis + expiryMillis}
	}
}

func (ms MemoryStore) GetData(key string) string {
	timeMillis := time.Now().UnixMilli()

	ms.DataMutex.RLock()
	defer ms.DataMutex.RUnlock()

	val, ok := ms.HashMap[key]
	if !ok {
		return ""
	}

	if val.Expiry != -1 && timeMillis > val.Expiry {
		delete(ms.HashMap, key)
		return ""
	}

	return val.Data
}
