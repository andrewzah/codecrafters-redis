package main

import (
	"sync"
	"time"
	"net"
)

type Role string

var (
	okResponse             = []byte("+OK\r\n")
	pongResponse           = []byte("+PONG\r\n")
	nullBulkStringResponse = []byte("$-1\r\n")
)

const (
	MasterRole Role = "master"
	ReplicaRole Role = "slave" // compatability

    okResponseStr string = "+OK\r\n"
)

type Store struct {
	Mutex sync.RWMutex
	Data  map[string]RedisValue
}

type RedisValue struct {
	Data string

	// in milliseconds
	Expiry int64
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

type InstanceMetadata struct {
	Role            Role
	ReplID          string
	ReplOffset      int

	ConnectedReplicasMutex sync.RWMutex
	ConnectedReplicas []chan<- string
}

func NewMetadata(args ServerArgs) *InstanceMetadata {
	if len(args.masterURL) == 0 {
		// master
		return &InstanceMetadata {
			MasterRole,
			RandStringBytes(40),
			0,
			sync.RWMutex{},
			[]chan<- string{},
		}
	} else {
		// replica
		return &InstanceMetadata {
			ReplicaRole,
			"?",
			-1,
			sync.RWMutex{},
			[]chan<- string{},
		}
	}
}

type ReplicaChan struct {
	Id string
	Conn net.Conn
}

type ReplicaSet struct {
	Mutex sync.RWMutex
	Replicas []ReplicaChan
}
