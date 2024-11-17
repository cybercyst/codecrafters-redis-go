package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type Command string

const (
	Ping Command = "ping"
	Echo Command = "echo"
	Set  Command = "set"
	Get  Command = "get"
	Info Command = "info"
)

type Server struct {
	address string
	port    int
}

type RedisServer struct {
	*Server
	store  Store
	master *Server
}

func (srv *RedisServer) Role() string {
	if srv.master != nil {
		return "slave"
	}

	return "master"
}

func NewRedisServer(address string, port int, master *Server) *RedisServer {
	return &RedisServer{
		Server: &Server{
			address: address,
			port:    port,
		},
		master: master,
		store:  *NewStore(),
	}
}

func (srv *RedisServer) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", srv.address, srv.port))
	if err != nil {
		return fmt.Errorf("failed to bind to port %d\n", srv.port)
	}
	fmt.Printf("Listening on port %d\n", srv.port)

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()

		go srv.handleClientConnection(conn)
	}
}

func parseMaster(replicaFlag string) *Server {
	if replicaFlag == "" {
		return nil
	}

	chunks := strings.SplitN(replicaFlag, " ", 2)

	address := chunks[0]
	portStr := chunks[1]

	if address == "" || portStr == "" {
		return nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil
	}

	return &Server{
		address: address,
		port:    port,
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("Exiting server...")
		close(c)
		os.Exit(0)
	}()

	portFlag := flag.Int("port", 6379, "the port for your redis server")
	replicaFlag := flag.String("replicaof", "", "the master host and master port")
	flag.Parse()

	master := parseMaster(*replicaFlag)

	srv := NewRedisServer("0.0.0.0", *portFlag, master)
	err := srv.Listen()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
