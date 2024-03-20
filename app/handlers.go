package main

import (
	"net"
	"strings"
)

func HandleEcho(cmd RedisCmd, conn net.Conn) error {
	response := encodeBulkString(strings.Join(cmd.Args, " "))

	return writeResponse([]byte(response), conn)
}

func HandlePing(conn net.Conn) error {
	return writeResponse(pongResponse, conn)
}

func HandleSet(cmd RedisCmd, conn net.Conn) error {
	memoryStore.DataMutex.Lock()
	defer memoryStore.DataMutex.Unlock()

	memoryStore.Data[cmd.Args[0]] = cmd.Args[1]

	return writeResponse(okResponse, conn)
}

func HandleGet(cmd RedisCmd, conn net.Conn) error {
	memoryStore.DataMutex.Lock()
	defer memoryStore.DataMutex.Unlock()

	if val, ok := memoryStore.Data[cmd.Args[0]]; ok {
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
