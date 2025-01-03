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
	"github.com/codecrafters-io/redis-starter-go/internal/replica"
)

var (
	replicationID = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
	emptyRbdFile  = "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="
)

func (srv *RedisServer) handlePing(conn net.Conn) error {
	return srv.WriteSimpleString(conn, "PONG")
}

func (srv *RedisServer) handleEcho(conn net.Conn, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("wrong number of arguments for 'echo' command")
	}
	return srv.WriteBulkString(conn, args[0])
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

	if srv.IsMaster() {
		for _, replica := range srv.replicas {
			msg := encoder.EncodeArrayBulkString([]string{"SET", key, val})
			err := replica.Write(msg)
			if err != nil {
				return fmt.Errorf("error writing to replica: %w", err)
			}
		}
	}

	if srv.IsSlave() {
		return nil
	}

	return srv.WriteSimpleString(conn, "OK")
}

func (srv *RedisServer) handleGet(conn net.Conn, args []string) error {
	key := args[0]
	return srv.WriteBulkString(conn, srv.store.Get(key))
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
		return srv.WriteBulkString(conn, resp)
	default:
		return fmt.Errorf("unknown info sub-command %s", subCmd)
	}
}

func (srv *RedisServer) handleReplConf(conn net.Conn, args []string) error {
	subCmd := args[0]
	switch subCmd {
	case "listening-port":
		portStr := args[1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("error parsing port %s: %w", portStr, err)
		}

		replica, err := replica.NewReplicaClient(conn, port)
		if err != nil {
			return fmt.Errorf("error creating replica client: %w", err)
		}
		srv.replicas = append(srv.replicas, replica)
	case "capa":
		break
	default:
		return fmt.Errorf("unknown repl conf sub-command %s", subCmd)
	}

	return srv.WriteSimpleString(conn, "OK")
}

func (srv *RedisServer) handlePsync(conn net.Conn) error {
	rbdFile, err := base64.StdEncoding.DecodeString(emptyRbdFile)
	if err != nil {
		return fmt.Errorf("error decoding empty RBD file: %w", err)
	}

	if err = srv.WriteSimpleString(conn, fmt.Sprintf("FULLRESYNC %s 0", replicationID)); err != nil {
		return fmt.Errorf("error sending FULLRESYNC: %w", err)
	}

	msg := []byte(fmt.Sprintf("$%d\r\n%s", len(rbdFile), rbdFile))
	return srv.Write(conn, msg)
}

func (srv *RedisServer) handle(conn net.Conn, parts []string) error {
	cmd := parts[0]
	args := parts[1:]
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
		return fmt.Errorf("unknown cmd %s", cmd)
	}
}

func (srv *RedisServer) WriteSimpleString(conn net.Conn, msg string) error {
	encodingMsg := encoder.EncodeSimpleString(msg)
	if err := srv.Write(conn, encodingMsg); err != nil {
		return err
	}
	return nil
}

func (srv *RedisServer) WriteBulkString(conn net.Conn, msg string) error {
	encodingMsg := encoder.EncodeBulkString(msg)
	if err := srv.Write(conn, encodingMsg); err != nil {
		return err
	}
	return nil
}

func (srv *RedisServer) Write(conn net.Conn, msg []byte) error {
	_, err := conn.Write(msg)
	return err
}

func (srv *RedisServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())

	for {
		parts, err := parseRESP(conn)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading from client: %s\n", err.Error())
			}
			return
		}

		if err := srv.handle(conn, parts); err != nil {
			fmt.Printf("Error handling command %s: %s\n", parts, err.Error())
		}
	}
}
