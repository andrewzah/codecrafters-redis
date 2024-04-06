package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/samber/lo"
)

func encodeSimpleString(input string) []byte {
	return []byte(fmt.Sprintf("$%s\r\n", input))
}

func encodeErrorString(input string) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", input))
}

func encodeArray(inputs []string) []byte {
	encodedInputs := lo.Map(inputs, func(s string, _idx int) string {
		return encodeBulkString(s)
	})

	return []byte(fmt.Sprintf("*%d\r\n%s", len(encodedInputs), strings.Join(encodedInputs, "")))
}

func encodeBulkString(input string) string {
	inputLength := utf8.RuneCountInString(input)

	return fmt.Sprintf("$%d\r\n%s\r\n", inputLength, input)
}
