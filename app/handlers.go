package main

import "fmt"

var memStore map[string]string = make(map[string]string)

func handlePing() string {
	return "+PONG\r\n"
}

func handleEcho(args []string) string {
	msg := args[0]
	return encodeBulkString(msg)
}

func handleSet(args []string) string {
	key := args[0]
	val := args[1]
	memStore[key] = val

	return encodeSimpleString("OK")
}

func handleGet(args []string) string {
	key := args[0]
	val := memStore[key]

	return encodeBulkString(val)
}

func handle(cmd string, args []string) (string, error) {
	switch Command(cmd) {
	case Ping:
		return handlePing(), nil
	case Echo:
		return handleEcho(args), nil
	case Set:
		return handleSet(args), nil
	case Get:
		return handleGet(args), nil
	default:
		return "", fmt.Errorf("unknown command %s", cmd)
	}
}
