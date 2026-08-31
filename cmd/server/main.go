package main

import (
	"context"
	"fmt"
	"github.com/Vupiad/mini-webserver-go/internal/tcp"
	"net"
)

type EchoHandler struct {
}

func (h *EchoHandler) ServeTCP(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}
	fmt.Printf("Received data: %s\n", string(data[:n]))

	_, err = conn.Write(data[:n])
	if err != nil {
		fmt.Printf("Error writing to connection: %v\n", err)
		return
	}

	fmt.Printf("Handling connection: %s\n", conn.RemoteAddr())
}

func main() {
	handler := &EchoHandler{}
	server := &tcp.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	server.Start(context.Background())
}
