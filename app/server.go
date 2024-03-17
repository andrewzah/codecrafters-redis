package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var (
	pongResponse = []byte("+PONG\r\n")
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	Debugf("Listening on 0.0.0.0:6379")
	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
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
			err = handlePing(conn)
		case "echo":
			err = handleEcho(redisCmd, conn)
		default:
			Errorf("Invalid redis command: [%s]", redisCmd.Name)
			return
		}

		if err != nil {
			Warningf("Error writing to client: ", err.Error())
			return
		}
	}

}

func handleEcho(cmd RedisCmd, conn net.Conn) error {
	response := encodeBulkString(strings.Join(cmd.Args, " "))

	return writeResponse([]byte(response), conn)
}

func handlePing(conn net.Conn) error {
	return writeResponse(pongResponse, conn)
}

func writeResponse(response []byte, conn net.Conn) error {
	_, err := conn.Write(response)

	if err != nil {
		return err
	} else {
		return nil
	}
}
