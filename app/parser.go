package main

import (
	"strings"

	"github.com/samber/lo"
)

type RedisCmd struct {
	Name string
	Args []string
}

// redis-cli delimits by whitespace
func parseRedisCmd(bytes []byte) RedisCmd {
	s := string(bytes[:])

	parsed := strings.Split(s, "\r\n")[1:]

	// filter out command length tokens like '$2'
	args := lo.Filter(parsed[2:], isArg)

	return RedisCmd{strings.ToLower(parsed[1]), args}
}

func isArg(str string, _idx int) bool {
	return !strings.HasPrefix(str, "$") && str != ""
}
