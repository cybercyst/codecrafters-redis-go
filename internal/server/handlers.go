package server

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/encoder"
)

func (srv *RedisServer) handlePing() string {
	return "+PONG\r\n"
}

func (srv *RedisServer) handleEcho(args []string) string {
	msg := args[0]
	return encoder.EncodeBulkString(msg)
}

func (srv *RedisServer) handleSet(args []string) (string, error) {
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
	return encoder.EncodeSimpleString("OK"), nil
}

func (srv *RedisServer) handleGet(args []string) string {
	key := args[0]
	return encoder.EncodeBulkString(srv.store.Get(key))
}

func (srv *RedisServer) handleInfo(args []string) (string, error) {
	subCmd := args[0]
	switch subCmd {
	case "replication":
		resp := strings.TrimSpace(fmt.Sprintf(`
role:%s
master_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb
master_repl_offset:0
		`, srv.Role()))
		return encoder.EncodeBulkString(resp), nil
	default:
		return "", fmt.Errorf("unknown info sub-command %s", subCmd)
	}
}

func (srv *RedisServer) handle(cmd string, args []string) (string, error) {
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

func (srv *RedisServer) handleClientConnection(conn net.Conn) {
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

		resp, err := srv.handle(cmd, args)
		if err != nil {
			fmt.Printf("Error handling command %s: %s\n", cmd, err.Error())
		}

		_, _ = conn.Write([]byte(resp))
	}
}
