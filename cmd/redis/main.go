package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/internal/replica"
	"github.com/codecrafters-io/redis-starter-go/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("error running redis", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	portFlag := flag.Int("port", 6379, "the port for your redis server")
	replicaFlag := flag.String("replicaof", "", "the master host and master port")
	flag.Parse()

	var masterClient *replica.MasterClient
	if *replicaFlag != "" {
		masterAddress, masterPort := parseReplica(*replicaFlag)
		m, err := replica.NewMasterClient(masterAddress, masterPort, *portFlag)
		if err != nil {
			return fmt.Errorf("error connecting to master: %w", err)
		}
		masterClient = m
		defer m.Close()
	}

	srv := server.NewRedisServer("0.0.0.0", *portFlag, masterClient)
	err := srv.Listen(ctx)
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	return nil
}

func parseReplica(replicaFlag string) (string, int) {
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
