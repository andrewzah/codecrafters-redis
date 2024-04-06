package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

func handshakePing(conn net.Conn) error {
	return nil
}

func handshakeReplConf(conn net.Conn, cmd string) error {
	return nil
}

func handleMasterHandshake(args ServerArgs) error {
	conn, err := net.Dial("tcp", args.masterURL)
	if err != nil {
		return fmt.Errorf("err connecting to master: %w", err)
	}

	err = handshakePing(conn)
	if err != nil {
		return fmt.Errorf("err with master handshake ping: %w", err)
	}

	request := encodeArray([]string{"ping"})
	Debugf("req: [%q]", request)

	_, err = conn.Write([]byte(request))
	if err != nil {
		return err
	}

	buf := make([]byte, 512)
	bytesRead, err := bufio.NewReader(conn).Read(buf)
	if err != nil {
		Fatalf("error receiving data from master")
	}

	if string(buf[:bytesRead]) != "+PONG\r\n" {
		msg := fmt.Sprintf("expected +PONG\\r\\n, received %q", buf[:bytesRead])
		return errors.New(msg)
	}

	return nil
}
