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

type ServerArgs struct {
	masterUrl string
	bindHost  string
	bindPort  uint
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

		a.masterUrl = fmt.Sprintf("%s:%d", *masterHost, masterPort)
	}
	return
}

func initMetadata(args ServerArgs) *InstanceMetadata {
	replID := ""
	role := NodeRole

	if len(args.masterUrl) == 0 {
		role = MasterRole
		replID = RandStringBytes(40)
	}

	metadata := InstanceMetadata{
		role,
		replID,
		0,
		0,
	}

	return &metadata
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

	metadata := initMetadata(args)
	store := Store{sync.RWMutex{}, map[string]RedisValue{}}

	if metadata.Role == NodeRole {
		go handlePingMaster(args.masterUrl)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn, &store, *metadata)
	}
}

func handleConnection(conn net.Conn, store *Store, metadata InstanceMetadata) {
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
			err = HandleInfo(redisCmd, conn, metadata)
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

			if e != nil {
				Fatalf("Error writing to client: %s", e.Error())
			}
		}
	}
}

func handlePingMaster(masterURL string) {
	conn, err := net.Dial("tcp", masterURL)
	if err != nil {
		Fatalf("frror connecting to master: %s", err.Error())
		os.Exit(1)
	}

	request := encodeArray([]string{"ping"})
	Debugf("req: [%q]", request)

	_, err = conn.Write([]byte(request))
	if err != nil {
		Fatalf("unable to write to master redis: %s", err.Error())
	}
}
