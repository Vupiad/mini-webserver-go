package tcp

import (
	"context"
	"net"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

type testHandler struct {
	handled bool
	mu      sync.Mutex
}

func (h *testHandler) ServeTCP(ctx context.Context, conn net.Conn) {
	h.mu.Lock()
	h.handled = true
	h.mu.Unlock()
	conn.Close()
}

func (h *testHandler) isHandled() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.handled
}

func TestServer(t *testing.T) {
	h := &testHandler{}
	s := &Server{
		Addr:    "127.0.0.1:18080", 
		Handler: h,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Start(context.Background())
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Connect to it
	conn, err := net.Dial("tcp", "127.0.0.1:18080")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	conn.Close()

	// Wait for handler to be called (it's asynchronous)
	time.Sleep(100 * time.Millisecond)

	if !h.isHandled() {
		t.Errorf("expected handler to be called")
	}

	// Trigger shutdown via signal
	p, err := os.FindProcess(os.Getpid())
	if err == nil {
		p.Signal(syscall.SIGINT)
	}

	// Wait for Start to return
	select {
	case <-errCh:
		// success
	case <-time.After(2 * time.Second):
		t.Errorf("timed out waiting for server to shut down")
	}
}
