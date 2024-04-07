package main

import (
	"net"
	"os"
	"bufio"
	"fmt"
	"errors"
	"bytes"
)

func handleHandshake(args ServerArgs, metadata *InstanceMetadata) {
	Debugf("started handshake")

	buf := make([]byte, 512)
	conn, err := net.Dial("tcp", args.masterURL)
	if err != nil {
		Fatalf("error connecting to master: %s", err.Error())
		os.Exit(1)
	}

	commands := [][][]byte{
		{ encodeArray([]string{"ping"}),
		  pongResponse },
		{ encodeArray([]string{"replconf", "listening-port", fmt.Sprint(args.bindPort)}),
		  okResponse },
		{ encodeArray([]string{"replconf", "capa", "psync2"}),
		  okResponse },
		{ encodeArray([]string{"psync", metadata.ReplID, fmt.Sprint(metadata.ReplOffset)}),
		  []byte{} },
	}

	for _, slice := range commands {
		Debugf("sending %s", slice[0])
		err = sendHandshakeCmd(conn, buf, slice[0], slice[1])

		if err != nil {
			Fatalf("error in setting up replica handshake: %s", err.Error())
			os.Exit(1)
		}
	}
}

func sendHandshakeCmd(conn net.Conn, buf []byte, request []byte, expectedResponse []byte) error {
	_, err := conn.Write([]byte(request))
	if err != nil {
		Fatalf("unable to write to master redis: %s", err.Error())
		return err
	}

	bytesRead, err := bufio.NewReader(conn).Read(buf)
	if err != nil {
		Fatalf("error receiving data from master: %s", err.Error())
		return err
	}

	if len(expectedResponse) > 0  {
		if !bytes.Equal(buf[:bytesRead], expectedResponse) {
			return errors.New(fmt.Sprintf("expected ok response, received %q", buf[:bytesRead]))
		}
	}

	return nil
}
