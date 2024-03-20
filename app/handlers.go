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

func HandleSet(cmd RedisCmd, conn net.Conn, store Store) error {
	switch len(cmd.Args) {
	case 1:
		return errors.New("Expected argument for SET command.")
	case 2:
		Debugf("inserting data")
		store.InsertData(cmd.Args[0], cmd.Args[1], -1)
	case 3:
		return errors.New("Expected argument for SET subcommand.")
	case 4:
		expiryMillis, err := strconv.ParseInt(cmd.Args[3], 10, 64)
		if err != nil {
			return errors.New("Unable to parse expiry argument into milliseconds (int64).")
		}
		store.InsertData(cmd.Args[0], cmd.Args[1], expiryMillis)
	default:
		return errors.New("Unexpected number of arguments for SET")
	}

	return writeResponse(okResponse, conn)
}

func HandleGet(cmd RedisCmd, conn net.Conn, store Store) error {
	val := store.GetData(cmd.Args[0])
	if val != "" {
		response := encodeBulkString(val)
		return writeResponse([]byte(response), conn)
	}
	return writeResponse(nullBulkStringResponse, conn)
}

func HandleInfo(cmd RedisCmd, conn net.Conn) error {
	if len(cmd.Args) < 1 {
		return errors.New("Expected subcommand for INFO command")
	}
	switch cmd.Args[0] {
	case "replication":
		response := encodeBulkString("role:master\nconnected_slaves:0")
		return writeResponse([]byte(response), conn)

	default:
		return errors.New("Unsupported subcommand for INFO command")
	}

}

func writeResponse(response []byte, conn net.Conn) error {
	_, err := conn.Write(response)

	if err != nil {
		return err
	} else {
		return nil
	}
}
