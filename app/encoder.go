package main

import (
	"fmt"
	"strings"
)

func encodeBulkString(val string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(val), val)
	return b.String()
}
