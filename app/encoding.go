package main

import (
	"fmt"
	"unicode/utf8"
)

func encodeString(input string) string {
	formatted := fmt.Sprintf("$%s\r\n", input)
	return formatted
}

func encodeBulkString(input string) string {
	inputLength := utf8.RuneCountInString(input)
	Debugf("inputLength: [%d], input: [%s]", inputLength, input)

	formatted := fmt.Sprintf("$%d\r\n%s\r\n", inputLength, input)
	return formatted
}
