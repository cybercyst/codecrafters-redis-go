package server

import (
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/encoder"
)

var (
	replicationID = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
	emptyRbdFile  = "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="
)

func (srv *RedisServer) handlePing(conn net.Conn) error {
	msg := encoder.EncodeSimpleString("PONG")
	_, err := conn.Write(msg)
	return err
}

func (srv *RedisServer) handleEcho(conn net.Conn, args []string) error {
	msg := encoder.EncodeBulkString(args[0])
	_, err := conn.Write(msg)
	return err
}

func (srv *RedisServer) handleSet(conn net.Conn, args []string) error {
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
			return err
		}
		expiry = time.Millisecond * time.Duration(expiryMs)
	}

	srv.store.Set(key, val, expiry)

	msg := encoder.EncodeSimpleString("OK")
	_, err := conn.Write(msg)
	return err
}

func (srv *RedisServer) handleGet(conn net.Conn, args []string) error {
	key := args[0]
	msg := encoder.EncodeBulkString(srv.store.Get(key))
	_, err := conn.Write(msg)
	return err
}

func (srv *RedisServer) handleInfo(conn net.Conn, args []string) error {
	subCmd := args[0]
	switch subCmd {
	case "replication":
		resp := strings.TrimSpace(fmt.Sprintf(`
role:%s
master_replid:%s
master_repl_offset:0
		`, srv.Role(), replicationID))
		msg := encoder.EncodeBulkString(resp)
		_, err := conn.Write(msg)
		return err
	default:
		return fmt.Errorf("unknown info sub-command %s", subCmd)
	}
}

func (srv *RedisServer) handleReplConf(conn net.Conn, args []string) error {
	msg := encoder.EncodeSimpleString("OK")
	_, err := conn.Write(msg)
	return err
}

func (srv *RedisServer) handlePsync(conn net.Conn) error {
	rbdFile, err := base64.StdEncoding.DecodeString(emptyRbdFile)
	if err != nil {
		return fmt.Errorf("error decoding empty RBD file: %w", err)
	}
	msg := encoder.EncodeSimpleString(fmt.Sprintf("FULLRESYNC %s 0", replicationID))
	_, err = conn.Write(msg)
	if err != nil {
		return fmt.Errorf("error sending FULLRESYNC: %w", err)
	}

	msg = []byte(fmt.Sprintf("$%d\r\n%s", len(rbdFile), rbdFile))
	_, err = conn.Write(msg)
	if err != nil {
		return fmt.Errorf("error sending rbd file: %w", err)
	}

	return nil
}

func (srv *RedisServer) handle(conn net.Conn, cmd string, args []string) error {
	switch Command(cmd) {
	case Ping:
		return srv.handlePing(conn)
	case Echo:
		return srv.handleEcho(conn, args)
	case Set:
		return srv.handleSet(conn, args)
	case Get:
		return srv.handleGet(conn, args)
	case Info:
		return srv.handleInfo(conn, args)
	case ReplConf:
		return srv.handleReplConf(conn, args)
	case PSync:
		return srv.handlePsync(conn)
	default:
		return fmt.Errorf("unknown command %s", cmd)
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

		if err := srv.handle(conn, cmd, args); err != nil {
			fmt.Printf("Error handling command %s: %s\n", cmd, err.Error())
		}
	}
}
