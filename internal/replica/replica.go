package replica

import "net"

func NewReplicaClient(conn net.Conn, listeningPort int) (*Client, error) {
	return &Client{
		Port: listeningPort,
		Conn: conn,
	}, nil
}
