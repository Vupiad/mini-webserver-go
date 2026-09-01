package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Handler interface {
	ServeTCP(ctx context.Context, conn net.Conn)
}

type Server struct {
	Addr    string
	Handler Handler
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.Addr, err)
	}
	defer listener.Close()

	fmt.Printf("Server is listening on %s\n", s.Addr)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	go func() {
		<-shutdownChan
		fmt.Println("Shutdown signal received, closing listener...")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			fmt.Printf("Failed to accept connection: %v\n", err)
			break
		}

		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			s.Handler.ServeTCP(ctx, c)
		}(conn)
	}

	wg.Wait()
	fmt.Println("All connections closed, server shutting down...")
	return nil
}
