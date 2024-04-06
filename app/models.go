package main

import (
	"sync"
	"time"
)

var (
	okResponse             = []byte("+OK\r\n")
    okResponseStr          = "+OK\r\n"
	pongResponse           = []byte("+PONG\r\n")
	nullBulkStringResponse = []byte("$-1\r\n")
)

const (
	MasterRole Role = "master"
	// for legacy compatibility purposes
	NodeRole Role = "slave"
)

type Role string

type Store struct {
	Mutex sync.RWMutex
	Data  map[string]RedisValue
}

type RedisValue struct {
	Data string

	// in milliseconds
	Expiry int64
}

type InstanceMetadata struct {
	Role            Role
	ReplID          string
	ReplOffset      uint
	ConnectedNodes uint
}

func (s *Store) InsertData(key, value string, expiryMillis int64) {
	timeMillis := time.Now().UnixMilli()

	Debugf("acquiring lock")
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	Debugf("actually inserting data")
	if expiryMillis == -1 {
		s.Data[key] = RedisValue{value, -1}
	} else {
		s.Data[key] = RedisValue{value, timeMillis + expiryMillis}
	}
}

func (s *Store) GetData(key string) string {
	timeMillis := time.Now().UnixMilli()

	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	val, ok := s.Data[key]
	if !ok {
		return ""
	}

	if val.Expiry != -1 && timeMillis > val.Expiry {
		delete(s.Data, key)
		return ""
	}

	return val.Data
}
