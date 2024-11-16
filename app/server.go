package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
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
	store   Store
}

func NewServer(address string, port int) *Server {
	return &Server{
		address: address,
		port:    port,
		store:   *NewStore(),
	}
}

func (srv *Server) Listen() error {
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
	flag.Parse()

	srv := NewServer("0.0.0.0", *portFlag)
	err := srv.Listen()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
