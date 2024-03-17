package main

import (
	"bufio"
	"fmt"
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

	buf := make([]byte, 512)
	_, err := bufio.NewReader(conn).Read(buf)

	if err != nil {
		return
	}

	Debugf("Conn from [%s]", conn.RemoteAddr().String())

	response := []byte("+PONG\r\n")
	conn.Write(response)
}
