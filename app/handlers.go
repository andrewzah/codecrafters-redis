package main

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

func HandleEcho(cmd RedisCmd, conn net.Conn) error {
	response := encodeBulkString(strings.Join(cmd.Args, " "))

	return writeResponse([]byte(response), conn)
}

func HandlePing(conn net.Conn) error {
	return writeResponse(pongResponse, conn)
}

func HandleSet(cmd RedisCmd, conn net.Conn, memoryStore MemoryStore) error {
	switch len(cmd.Args) {
	case 1:
		return errors.New("Expected argument for SET command.")
	case 2:
		Debugf("inserting data")
		memoryStore.InsertData(cmd.Args[0], cmd.Args[1], -1)
	case 3:
		return errors.New("Expected argument for SET subcommand.")
	case 4:
		expiryMillis, err := strconv.ParseInt(cmd.Args[3], 10, 64)
		if err != nil {
			return errors.New("Unable to parse expiry argument into milliseconds (int64).")
		}
		memoryStore.InsertData(cmd.Args[0], cmd.Args[1], expiryMillis)
	default:
		return errors.New("Unexpected number of arguments for SET")
	}

	return writeResponse(okResponse, conn)
}

func HandleGet(cmd RedisCmd, conn net.Conn, memoryStore MemoryStore) error {
	val := memoryStore.GetData(cmd.Args[0])
	if val != "" {
		response := encodeBulkString(val)
		return writeResponse([]byte(response), conn)
	}
	return writeResponse(nullBulkStringResponse, conn)
}

func writeResponse(response []byte, conn net.Conn) error {
	_, err := conn.Write(response)

	if err != nil {
		return err
	} else {
		return nil
	}
}
