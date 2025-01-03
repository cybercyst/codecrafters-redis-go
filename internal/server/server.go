package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/replica"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

type RedisServer struct {
	address      string
	port         int
	store        store.Store
	masterClient *replica.Client
	replicas     []*replica.Client
}

func (srv *RedisServer) Role() string {
	if srv.masterClient != nil {
		return "slave"
	}

	return "master"
}

func (srv *RedisServer) IsStandalone() bool {
	return srv.IsMaster() && srv.replicas == nil
}

func (srv *RedisServer) IsMaster() bool {
	return srv.Role() == "master"
}

func (srv *RedisServer) IsSlave() bool {
	return srv.Role() == "slave"
}

func NewRedisServer(address string, port int, masterClient *replica.Client) *RedisServer {
	return &RedisServer{
		address:      address,
		port:         port,
		masterClient: masterClient,
		store:        *store.NewStore(),
		replicas:     []*replica.Client{},
	}
}

func (srv *RedisServer) Listen(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", srv.address, srv.port))
	if err != nil {
		slog.Error("failed to bind to port", slog.Int("port", srv.port))
		return err
	}
	defer listener.Close()
	slog.Info("Redis started", slog.Int("port", srv.port))

	go func() {
		<-ctx.Done()
		fmt.Println("\nShutting down server...")
		listener.Close()
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("failed to receive data", slog.Any("error", err))
			return err
		}
		go srv.handleConnection(conn)
	}
}
