package replica

import (
	"fmt"
	"net"
)

type MasterClient struct {
	address string
	port    int
	conn    net.Conn
}

func (m *MasterClient) Close() {
	m.conn.Close()
}

func NewMasterClient(address string, port int) (*MasterClient, error) {
	if address == "" {
		return nil, nil
	}
	if port == 0 {
		return nil, nil
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, fmt.Errorf("error establishing connection to %s:%d: %w", address, port, err)
	}

	_, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		return nil, fmt.Errorf("error sending PING: %w", err)
	}

	return &MasterClient{
		address: address,
		port:    port,
		conn:    conn,
	}, nil
}
