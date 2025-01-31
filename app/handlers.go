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

func HandleGet(cmd RedisCmd, conn net.Conn, store *Store) error {
	val := store.GetData(cmd.Args[0])
	if val != "" {
		response := encodeBulkString(val)
		return writeResponse([]byte(response), conn)
	}
	return writeResponse(nullBulkStringResponse, conn)
}

func HandleInfo(cmd RedisCmd, conn net.Conn, md *InstanceMetadata) error {
	if len(cmd.Args) < 1 {
		return errors.New("expected subcommand for INFO command")
	}

	md.ConnectedReplicasMutex.RLock()
	defer md.ConnectedReplicasMutex.RUnlock()
    
	switch cmd.Args[0] {
	case "replication":
		response := fmt.Sprintf("role:%s\nconnected_slaves:%d\nmaster_replid:%s\nmaster_repl_offset:%d",
			md.Role, len(md.ConnectedReplicas),
			md.ReplID, md.ReplOffset)

		encodedResponse := encodeBulkString(response)
		return writeResponse([]byte(encodedResponse), conn)

	default:
		return errors.New("unsupported subcommand for INFO command")
	}
}


func HandleSet(cmd RedisCmd, conn net.Conn, store *Store) error {
	switch len(cmd.Args) {
	case 1:
		return errors.New("expected argument for SET command")
	case 2:
		Debugf("inserting data")
		store.InsertData(cmd.Args[0], cmd.Args[1], -1)
	case 3:
		return errors.New("expected argument for SET subcommand")
	case 4:
		expiryMillis, err := strconv.ParseInt(cmd.Args[3], 10, 64)
		if err != nil {
			return errors.New("unable to parse expiry argument into milliseconds (int64)")
		}
		store.InsertData(cmd.Args[0], cmd.Args[1], expiryMillis)
	default:
		return errors.New("unexpected number of arguments for SET")
	}

	return writeResponse(okResponse, conn)
}

func HandleReplconf(conn net.Conn) error {
	return writeResponse(okResponse, conn)
}

func HandlePing(conn net.Conn) error {
	return writeResponse(pongResponse, conn)
}

func HandlePsync(conn net.Conn, md *InstanceMetadata) error {
    response := fmt.Sprintf("FULLRESYNC %s 0", md.ReplID)

    err :=  writeResponse(encodeSimpleString(response), conn)
    if err != nil {
        return err
    }

    return writeResponse(emptyRDBFileResponse(), conn)
}

func writeResponse(response []byte, conn net.Conn) error {
	_, err := conn.Write(response)
	return err
}
