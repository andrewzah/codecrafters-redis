package main

import (
	"fmt"
	"unicode/utf8"
)

func encodeSimpleString(input string) string {
	formatted := fmt.Sprintf("$%s\r\n", input)
	return formatted
}

func encodeErrorString(input string) string {
	formatted := fmt.Sprintf("-%s\r\n", input)
	return formatted
}

func encodeBulkString(input string) string {
	inputLength := utf8.RuneCountInString(input)

	formatted := fmt.Sprintf("$%d\r\n%s\r\n", inputLength, input)
	return formatted
}
