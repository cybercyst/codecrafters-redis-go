package encoder

import (
	"fmt"
)

func EncodeArrayBulkString(val []string) []byte {
	resp := []byte(fmt.Sprintf("*%d\r\n", len(val)))

	for _, v := range val {
		resp = append(resp, EncodeBulkString(v)...)
	}

	return []byte(resp)
}

func EncodeBulkString(val string) []byte {
	if val == "" {
		return []byte("$-1\r\n")
	}
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val))
}

func EncodeSimpleString(val string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", val))
}
