package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
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

	var masterClient *replica.Client
	if *replicaFlag != "" {
		m, err := replica.NewMasterClient(ctx, *replicaFlag, *portFlag)
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
