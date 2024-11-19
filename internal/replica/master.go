package replica

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/encoder"
)

type MasterClient struct {
	address string
	port    int
	conn    net.Conn
}

func (m *MasterClient) Close() {
	m.conn.Close()
}

func (m *MasterClient) SendMessageAndGetResponse(req string) (string, error) {
	encodedReq := encoder.EncodeArrayBulkString(strings.Split(req, " "))
	_, err := (m.conn).Write(encodedReq)
	if err != nil {
		return "", fmt.Errorf("error sending message %s: %w", req, err)
	}

	message, err := bufio.NewReader(m.conn).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error getting reply: %w", err)
	}
	return message, nil
}

func (m *MasterClient) SendMessageAndExpectResponse(req, resp string) error {
	encodedReq := encoder.EncodeArrayBulkString(strings.Split(req, " "))
	_, err := (m.conn).Write(encodedReq)
	if err != nil {
		return fmt.Errorf("error sending message %s: %w", req, err)
	}

	expectedResp := encoder.EncodeSimpleString(resp)
	message, _ := bufio.NewReader(m.conn).ReadString('\n')
	if string(expectedResp) != message {
		return fmt.Errorf("expected response %s, got %s", resp, message)
	}

	return nil
}

func NewMasterClient(masterAddress string, masterPort, listeningPort int) (*MasterClient, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", masterAddress, masterPort))
	if err != nil {
		return nil, fmt.Errorf("error establishing connection to %s:%d: %w", masterAddress, masterPort, err)
	}

	client := &MasterClient{
		address: masterAddress,
		port:    masterPort,
		conn:    conn,
	}

	if err = client.SendMessageAndExpectResponse("PING", "PONG"); err != nil {
		return nil, fmt.Errorf("error sending PING: %w", err)
	}

	msg := fmt.Sprintf("REPLCONF listening-port %d", listeningPort)
	if err = client.SendMessageAndExpectResponse(msg, "OK"); err != nil {
		return nil, fmt.Errorf("error sending %s: %w", msg, err)
	}

	msg = "REPLCONF capa psync2"
	if err = client.SendMessageAndExpectResponse(msg, "OK"); err != nil {
		return nil, fmt.Errorf("error sending %s: %w", msg, err)
	}

	resp, err := client.SendMessageAndGetResponse("PSYNC ? -1")
	if err != nil {
		return nil, fmt.Errorf("error sending PSYNC: %w", err)
	}
	slog.Info("PSYNC response", slog.String("resp", resp))

	return client, nil
}
