package main

import (
	"fmt"
)

func encodeBulkString(val string) string {
	if val == "" {
		return "$-1\r\n"
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
}

func encodeSimpleString(val string) string {
	return fmt.Sprintf("+%s\r\n", val)
}
