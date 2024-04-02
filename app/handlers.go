package main

import (
	"errors"
	"fmt"
	"io"
	"net"
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
	case "px":
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

func (srv *Server) handleInfo(args []string) (string, error) {
	subCmd := args[0]
	switch subCmd {
	case "replication":
		return encodeBulkString("role:master"), nil
	default:
		return "", errors.New(fmt.Sprintf("unknown info sub-command %s", subCmd))
	}
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
	case Info:
		return srv.handleInfo(args)
	default:
		return "", fmt.Errorf("unknown command %s", cmd)
	}
}

func (srv *Server) handleClientConnection(conn net.Conn) {
	defer conn.Close()

	for {
		cmd, args, err := parseRequest(conn)
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("Error reading from client: %s\n", err.Error())
			return
		}

		// fmt.Println("Received from client:")
		// fmt.Println("Cmd: ", cmd)
		// fmt.Println("Args: ", args)

		resp, err := srv.handle(cmd, args)
		if err != nil {
			fmt.Printf("Error handling command %s: %s\n", cmd, err.Error())
		}

		_, _ = conn.Write([]byte(resp))
	}
}
