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

func handleClient(conn net.Conn) {
	buff := make([]byte, 1024)
	for {
		_, err := conn.Read(buff)
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading from client: ", err.Error())
			os.Exit(1)
		}
		// fmt.Println("Received message from client: ")
		// fmt.Println(string(buff))

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error writing to client: ", err.Error())
			os.Exit(1)
		}
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

		go handleClient(conn)
	}
}
