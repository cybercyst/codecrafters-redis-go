package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const redisPort = 6379

type Command string

const (
	Ping Command = "ping"
	Echo         = "echo"
	Set          = "set"
	Get          = "get"
)

func handleClientConnection(conn net.Conn) {
	defer conn.Close()

	for {
		cmd, args, err := parseRequest(conn)
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("Error reading from client: %s\n", err.Error())
			return
		}

		fmt.Println(cmd)
		fmt.Println(args)

		resp, err := handle(cmd, args)
		if err != nil {
			fmt.Printf("Error handling command %s: %s\n", cmd, err.Error())
		}

		_, _ = conn.Write([]byte(resp))
	}
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("Exiting server...")
		close(c)
		os.Exit(0)
	}()

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", redisPort))
	if err != nil {
		fmt.Printf("Failed to bind to port %d\n", redisPort)
		os.Exit(1)
	}
	fmt.Printf("Listening on port %d\n", redisPort)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		defer conn.Close()

		go handleClientConnection(conn)
	}
}
