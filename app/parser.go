package main

import (
	"strings"

	"github.com/samber/lo"
)

type RedisCmd struct {
	Name string
	Args []string
}

func parseRedisCmd(bytes []byte) RedisCmd {
	s := string(bytes[:])
	parsed := strings.Split(s, "\r\n")[1:]

	args := lo.Filter(parsed[2:], func(str string, _idx int) bool {
		return !strings.HasPrefix(str, "$") && str != ""
	})

	return RedisCmd{strings.ToLower(parsed[1]), args}
}
