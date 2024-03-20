package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

var (
	okResponse             = []byte("+OK\r\n")
	pongResponse           = []byte("+PONG\r\n")
	nullBulkStringResponse = []byte("$-1\r\n")
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	memoryStore := MemoryStore{map[string]RedisValue{}, sync.RWMutex{}}

	Debugf("Listening on 0.0.0.0:6379")
	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn, memoryStore)
	}
}

func handleConnection(conn net.Conn, memoryStore MemoryStore) {
	defer conn.Close()
	buf := make([]byte, 512)

	for {
		bytes_read, err := bufio.NewReader(conn).Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			Warningf("Error reading buffer: ", err.Error())
			return
		}

		redisCmd := parseRedisCmd(buf[:bytes_read])

		switch redisCmd.Name {
		case "ping":
			err = HandlePing(conn)
		case "echo":
			err = HandleEcho(redisCmd, conn)
		case "set":
			err = HandleSet(redisCmd, conn, memoryStore)
		case "get":
			err = HandleGet(redisCmd, conn, memoryStore)
		default:
			Errorf("Invalid redis command: [%s]", redisCmd.Name)
			return
		}

		if err != nil {
			Errorf("Error writing to client: ", err.Error())
			return
		}
	}

}
