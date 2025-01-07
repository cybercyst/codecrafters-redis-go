package replica

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
)

func NewMasterClient(ctx context.Context, replicaFlag string, listeningPort int) (*Client, error) {
	masterAddress, masterPort := parseReplicaFlag(replicaFlag)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", masterAddress, masterPort))
	if err != nil {
		return nil, fmt.Errorf("error establishing connection to %s:%d: %w", masterAddress, masterPort, err)
	}

	client := &Client{
		Address: masterAddress,
		Port:    masterPort,
		Conn:    conn,
	}

	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down master connection...")
		client.Close()
	}()

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

	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(message)

	return client, nil
}

func parseReplicaFlag(replicaFlag string) (string, int) {
	if replicaFlag == "" {
		return "", 0
	}

	chunks := strings.SplitN(replicaFlag, " ", 2)

	address := chunks[0]
	portStr := chunks[1]

	if address == "" || portStr == "" {
		return "", 0
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0
	}

	return address, port
}
