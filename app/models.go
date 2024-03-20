package main

import (
	"sync"
	"time"
)

type Store struct {
	Map      map[string]RedisValue
	MapMutex sync.RWMutex
}

type RedisValue struct {
	Data string

	// in milliseconds
	Expiry int64
}

func (s Store) InsertData(key, value string, expiryMillis int64) {
	timeMillis := time.Now().UnixMilli()

	Debugf("acquiring lock")
	s.MapMutex.Lock()
	defer s.MapMutex.Unlock()

	Debugf("actually inserting data")
	if expiryMillis == -1 {
		s.Map[key] = RedisValue{value, -1}
	} else {
		s.Map[key] = RedisValue{value, timeMillis + expiryMillis}
	}
}

func (s Store) GetData(key string) string {
	timeMillis := time.Now().UnixMilli()

	s.MapMutex.RLock()
	defer s.MapMutex.RUnlock()

	val, ok := s.Map[key]
	if !ok {
		return ""
	}

	if val.Expiry != -1 && timeMillis > val.Expiry {
		delete(s.Map, key)
		return ""
	}

	return val.Data
}

type InfoReplication struct {
	role string
}
