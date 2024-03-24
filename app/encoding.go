package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/samber/lo"
)

func encodeSimpleString(input string) string {
	formatted := fmt.Sprintf("$%s\r\n", input)
	return formatted
}

func encodeErrorString(input string) string {
	formatted := fmt.Sprintf("-%s\r\n", input)
	return formatted
}

func encodeArray(inputs []string) string {
	encodedInputs := lo.Map(inputs, func(s string, _idx int) string {
		return encodeBulkString(s)
	})

	formatted := fmt.Sprintf("*%d\r\n%s", len(encodedInputs), strings.Join(encodedInputs, ""))
	return formatted
}

func encodeBulkString(input string) string {
	inputLength := utf8.RuneCountInString(input)

	formatted := fmt.Sprintf("$%d\r\n%s\r\n", inputLength, input)
	return formatted
}
