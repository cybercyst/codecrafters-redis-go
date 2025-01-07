package replica

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/encoder"
)

type Client struct {
	Address string
	Port    int
	Conn    net.Conn
}

func (m *Client) Write(msg []byte) error {
	_, err := m.Conn.Write(msg)
	return err
}

func (m *Client) Close() {
	m.Conn.Close()
}

func (m *Client) SendMessageAndGetResponse(req string) (string, error) {
	encodedReq := encoder.EncodeArrayBulkString(strings.Split(req, " "))
	_, err := (m.Conn).Write(encodedReq)
	if err != nil {
		return "", fmt.Errorf("error sending message %s: %w", req, err)
	}

	message, err := bufio.NewReader(m.Conn).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error getting reply: %w", err)
	}
	return message, nil
}

func (m *Client) SendMessageAndExpectResponse(req, resp string) error {
	encodedReq := encoder.EncodeArrayBulkString(strings.Split(req, " "))
	_, err := (m.Conn).Write(encodedReq)
	if err != nil {
		return fmt.Errorf("error sending message %s: %w", req, err)
	}

	expectedResp := encoder.EncodeSimpleString(resp)
	message, _ := bufio.NewReader(m.Conn).ReadString('\n')
	if string(expectedResp) != message {
		return fmt.Errorf("expected response %s, got %s", resp, message)
	}

	return nil
}
