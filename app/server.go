package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const redisPort = 6379

type Command string

const (
	Ping Command = "ping"
	Echo         = "echo"
)

func handlePing(conn net.Conn) {
	_, err := conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		fmt.Printf("Error writing to client: %s\n", err.Error())
	}
}

func handleEcho(conn net.Conn, msg string) {
	_, err := conn.Write([]byte(encodeBulkString(msg)))
	if err != nil {
		fmt.Printf("Error writing to client: %s\n", err.Error())
	}
}

func handleClientConnection(conn net.Conn) {
	defer conn.Close()

	for {
		cmds, err := parseRequest(conn)
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("Error reading from client: %s\n", err.Error())
			return
		}

		if len(cmds) < 1 {
			fmt.Println("No commands parsed!")
			return
		}

		// fmt.Println("Received message from client: ")
		// fmt.Println(cmds)

		cmdStr, ok := cmds[0].(string)
		if !ok {
			fmt.Printf("unable to parse command as string: %s", cmdStr)
			return
		}
		switch Command(strings.ToLower(cmdStr)) {
		case Ping:
			handlePing(conn)
			break
		case Echo:
			var b strings.Builder
			for _, cmd := range cmds[1:] {
				cmdStr, ok := cmd.(string)
				if !ok {
					fmt.Printf("error parsing command as string: %v", cmd)
					break
				}
				fmt.Fprintf(&b, "%s ", cmdStr)
			}
			handleEcho(conn, strings.TrimSpace(b.String()))
			break
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

		go handleClientConnection(conn)
	}
}
