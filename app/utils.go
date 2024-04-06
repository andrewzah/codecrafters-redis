package main

import (
	"crypto/rand"
	"fmt"
	"os"
    "encoding/hex"
)

const letters = "abcdefghijklmnopqrstuvwxyz123456789"

func RandStringBytes(l int) string {
	b := make([]byte, l)

	_, err := rand.Read(b)
	if err != nil {
		panic("Could not generate id for master")
	}

	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}

	return string(b)
}

// eventually replace with real RDB functionality
func emptyRDBFileResponse() []byte {
    hexString := "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"
    bytes, _ := hex.DecodeString(hexString)
    response := fmt.Sprintf("$%d\r\n%s", len(bytes), bytes)
    return []byte(response)
}

// logs at level
func LogAtLevel(level, format string, a ...interface{}) {
	msg := fmt.Sprintf("[%s] %s\n", level, format)
	fmt.Fprintf(os.Stdout, msg, a...)
}

func Fatalf(format string, a ...interface{})   { LogAtLevel("FTL", format, a...) }
func Warningf(format string, a ...interface{}) { LogAtLevel("WRN", format, a...) }
func Debugf(format string, a ...interface{})   { LogAtLevel("DBG", format, a...) }
func Errorf(format string, a ...interface{})   { LogAtLevel("ERR", format, a...) }
func Infof(format string, a ...interface{})    { LogAtLevel("INF", format, a...) }
