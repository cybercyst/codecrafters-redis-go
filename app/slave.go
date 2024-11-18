package main

import (
	"fmt"
	"net"
)

type Slave struct {
	*Server
	conn net.Conn
}

func NewSlave(server *Server) (*Slave, error) {
	if server == nil {
		return nil, nil
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.address, server.port))
	if err != nil {
		return nil, fmt.Errorf("error establishing connection to %s:%d: %w", server.address, server.port, err)
	}

	_, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		return nil, fmt.Errorf("error sending PING: %w", err)
	}

	return &Slave{
		Server: server,
		conn:   conn,
	}, nil
}
