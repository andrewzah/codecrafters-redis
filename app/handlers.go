package main

import (
	"errors"
	"fmt"
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

func HandleSet(cmd RedisCmd, conn net.Conn, ctx *AppContext) error {
	switch len(cmd.Args) {
	case 1:
		return errors.New("expected argument for SET command")
	case 2:
		Debugf("inserting data")
		ctx.InsertData(cmd.Args[0], cmd.Args[1], -1)
	case 3:
		return errors.New("expected argument for SET subcommand")
	case 4:
		expiryMillis, err := strconv.ParseInt(cmd.Args[3], 10, 64)
		if err != nil {
			return errors.New("unable to parse expiry argument into milliseconds (int64)")
		}
		ctx.InsertData(cmd.Args[0], cmd.Args[1], expiryMillis)
	default:
		return errors.New("unexpected number of arguments for SET")
	}

	return writeResponse(okResponse, conn)
}

func HandleGet(cmd RedisCmd, conn net.Conn, ctx *AppContext) error {
	val := ctx.GetData(cmd.Args[0])
	if val != "" {
		response := encodeBulkString(val)
		return writeResponse([]byte(response), conn)
	}
	return writeResponse(nullBulkStringResponse, conn)
}

func HandleInfo(cmd RedisCmd, conn net.Conn, ctx *AppContext) error {
	if len(cmd.Args) < 1 {
		return errors.New("expected subcommand for INFO command")
	}
	switch cmd.Args[0] {
	case "replication":
		response := fmt.Sprintf("role:%s\nconnected_slaves:%d\nmaster_replid:%s\nmaster_repl_offset:%d",
			ctx.Metadata.Role, ctx.Metadata.ConnectedSlaves,
			ctx.Metadata.ReplID, ctx.Metadata.ReplOffset)

		encodedResponse := encodeBulkString(response)
		return writeResponse([]byte(encodedResponse), conn)

	default:
		return errors.New("unsupported subcommand for INFO command")
	}
}

func writeResponse(response []byte, conn net.Conn) error {
	_, err := conn.Write(response)
	return err
}
