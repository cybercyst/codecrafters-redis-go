package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/replica"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

type RedisServer struct {
	address      string
	port         int
	store        store.Store
	masterClient *replica.MasterClient
}

func (srv *RedisServer) Role() string {
	if srv.masterClient != nil {
		return "slave"
	}

	return "master"
}

func NewRedisServer(address string, port int, masterClient *replica.MasterClient) *RedisServer {
	return &RedisServer{
		address:      address,
		port:         port,
		masterClient: masterClient,
		store:        *store.NewStore(),
	}
}

func (srv *RedisServer) Listen(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", srv.address, srv.port))
		if err != nil {
			slog.Error("failed to bind to port", slog.Int("port", srv.port))
			cancel()
		}
		defer listener.Close()
		slog.Info("Redis started", slog.Int("port", srv.port))

		for {
			conn, err := listener.Accept()
			if err != nil {
				slog.Error("failed to receive data", slog.Any("error", err))
				cancel()
			}
			go srv.handleClientConnection(conn)
		}
	}()

	<-ctx.Done()

	return nil
}
