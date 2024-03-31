package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (srv *Server) handlePing() string {
	return "+PONG\r\n"
}

func (srv *Server) handleEcho(args []string) string {
	msg := args[0]
	return encodeBulkString(msg)
}

func (srv *Server) handleSet(args []string) (string, error) {
	key := args[0]
	val := args[1]

	subCmd := ""
	if len(args) > 2 {
		subCmd = strings.ToLower(args[2])
	}

	var expiry time.Duration
	switch subCmd {
	case "ex":
		expiryMs, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return "", err
		}
		expiry = time.Millisecond * time.Duration(expiryMs)
	}

	srv.store.Set(key, val, expiry)
	return encodeSimpleString("OK"), nil
}

func (srv *Server) handleGet(args []string) string {
	key := args[0]
	return encodeBulkString(srv.store.Get(key))
}

func (srv *Server) handle(cmd string, args []string) (string, error) {
	switch Command(cmd) {
	case Ping:
		return srv.handlePing(), nil
	case Echo:
		return srv.handleEcho(args), nil
	case Set:
		return srv.handleSet(args)
	case Get:
		return srv.handleGet(args), nil
	default:
		return "", fmt.Errorf("unknown command %s", cmd)
	}
}
