package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
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
	pongResponse := []byte("+PONG\r\n")
	buf := make([]byte, 512)

	for {
		Debugf("[%s]: %s", conn.RemoteAddr().String(), buf)

		_, err := bufio.NewReader(conn).Read(buf)

		if err == io.EOF {
			break
		}

		if err != nil {
			Warningf("Error reading buffer: ", err.Error())
			return
		}

		_, err = conn.Write(pongResponse)
		if err != nil {
			Warningf("Error writing to client: ", err.Error())
			return
		}
	}

}
