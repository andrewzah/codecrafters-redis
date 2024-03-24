package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
)

var (
	okResponse             = []byte("+OK\r\n")
	pongResponse           = []byte("+PONG\r\n")
	nullBulkStringResponse = []byte("$-1\r\n")
)

type ServerArgs struct {
	master   string
	bindHost string
	bindPort uint
}

func parseArgs() (a ServerArgs, e error) {
	bindHost := flag.String("bind-host", "0.0.0.0", "host to bind on")
	port := flag.Uint("port", 6379, "port for server")
	replicaOf := flag.String("replicaof", "", "the host and port of the master server")

	flag.Parse()

	a.bindHost = *bindHost
	a.bindPort = *port

	if len(*replicaOf) > 0 {
		Debugf("replicaof: %q", *replicaOf)
		masterHost := replicaOf
		masterPortStr := flag.Args()[0]

		if len(masterPortStr) == 0 {
			e = errors.New("replica string must contain an ip and port separated by a space")
			return
		}

		masterPort, err := strconv.Atoi(masterPortStr)
		if err != nil {
			e = err
			return
		}

		a.master = fmt.Sprintf("%s:%d", *masterHost, masterPort)
	}
	return
}

func createContext(args ServerArgs) *AppContext {
	replID := ""
	role := "slave"

	if len(args.master) == 0 {
		role = "master"
		replID = RandStringBytes(40)
	}

	metadata := AppMetadata{
		role,
		replID,
		0,
		0,
	}

	ctx := AppContext{
		sync.RWMutex{},
		map[string]RedisValue{},
		sync.RWMutex{},
		metadata,
	}

	return &ctx
}

func main() {
	args, err := parseArgs()
	if err != nil {
		fmt.Printf("Failed to parse args:\n  %s", err.Error())
	}

	bind := fmt.Sprintf("%s:%d", args.bindHost, args.bindPort)
	l, err := net.Listen("tcp", bind)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	Infof("Listening on %s", bind)

	ctx := createContext(args)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn, ctx)
	}
}

func handleConnection(conn net.Conn, ctx *AppContext) {
	defer conn.Close()
	buf := make([]byte, 512)

	for {
		bytesRead, err := bufio.NewReader(conn).Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			Warningf("Error reading buffer: ", err.Error())
			return
		}

		redisCmd := parseRedisCmd(buf[:bytesRead])

		switch redisCmd.Name {
		case "ping":
			err = HandlePing(conn)
		case "echo":
			err = HandleEcho(redisCmd, conn)
		case "info":
			err = HandleInfo(redisCmd, conn, ctx)
		case "set":
			err = HandleSet(redisCmd, conn, ctx)
		case "get":
			err = HandleGet(redisCmd, conn, ctx)
		default:
			Errorf("Invalid redis command: [%s]", redisCmd.Name)
			return
		}

		if err != nil {
			Errorf(err.Error())
			_, e := conn.Write([]byte(encodeErrorString(err.Error())))

			if e != nil {
				Fatalf("Error writing to client: %s", e.Error())
			}
		}
	}
}
