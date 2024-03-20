package main

import (
	"bufio"
	"flag"
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
	port := flag.String("port", "6379", "port for server")
	bindHost := flag.String("bind-host", "0.0.0.0", "host to bind on")
	flag.Parse()

	formattedHost := fmt.Sprintf("%s:%s", *bindHost, *port)
	l, err := net.Listen("tcp", formattedHost)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	Infof("Listening on %s", formattedHost)
	store := Store{map[string]RedisValue{}, sync.RWMutex{}}

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn, store)
	}
}

func handleConnection(conn net.Conn, store Store) {
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
		case "info":
			err = HandleInfo(redisCmd, conn)
		case "set":
			err = HandleSet(redisCmd, conn, store)
		case "get":
			err = HandleGet(redisCmd, conn, store)
		default:
			Errorf("Invalid redis command: [%s]", redisCmd.Name)
			return
		}

		if err != nil {
			Errorf(err.Error())
			_, e := conn.Write([]byte(encodeErrorString(err.Error())))
			if err != nil {
				Fatalf("Error writing to client: %s", e.Error())
			}
		}
	}
}
