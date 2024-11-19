package server

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type TokenType byte

const (
	BulkString TokenType = '$'
	Array      TokenType = '*'
	Integer    TokenType = ':'
)

func parseRequest(reader io.Reader) (string, []string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	tokens, err := parseToken(scanner)
	if err == io.EOF {
		return "", nil, err
	}

	if err != nil {
		return "", nil, fmt.Errorf("error parsing request: %v", err)
	}

	switch tokens.(type) {
	case []any:
		tokens := tokens.([]any)
		cmd := strings.ToLower(tokens[0].(string))
		args := []string{}
		if len(tokens) > 1 {
			for _, token := range tokens[1:] {
				args = append(args, token.(string))
			}
		}
		return cmd, args, nil
	default:
		return "", nil, fmt.Errorf("failed to parse slice of commands")
	}
}

func parseToken(scanner *bufio.Scanner) (any, error) {
	hasMore := scanner.Scan()
	if !hasMore {
		return nil, io.EOF
	}

	token := scanner.Text()
	if token == "" {
		return nil, fmt.Errorf("no token parsed")
	}

	switch TokenType(token[0]) {
	case Array:
		return parseArray(scanner)
	case BulkString:
		return parseBulkString(scanner)
	}
	return nil, fmt.Errorf("unknown token: %s", token)
}

func parseArray(scanner *bufio.Scanner) ([]any, error) {
	token := scanner.Text()
	if TokenType(token[0]) != Array {
		return nil, fmt.Errorf("token %s not an array", token)
	}

	arraySize, err := strconv.ParseInt(token[1:], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing array size from token: %v", err)
	}

	arr := make([]any, arraySize)
	for i := int64(0); i < arraySize; i++ {
		arr[i], err = parseToken(scanner)
		if err != nil {
			return nil, err
		}
	}

	return arr, nil
}

func parseBulkString(scanner *bufio.Scanner) (string, error) {
	token := scanner.Text()
	if TokenType(token[0]) != BulkString {
		return "", fmt.Errorf("token %s is not a bulk string", token)
	}

	stringSize, err := strconv.ParseInt(token[1:], 10, 64)
	if err != nil {
		return "", fmt.Errorf("error parsing string size from token: %v", err)
	}

	scanner.Scan()
	stringVal := scanner.Text()
	return stringVal[:stringSize], nil
}
