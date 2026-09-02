package http

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"
)

// mockConn implements net.Conn for testing
type mockConn struct {
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
	closed      bool
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuffer.Read(b)
}
func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuffer.Write(b)
}
func (m *mockConn) Close() error { 
	m.closed = true
	return nil 
}
func (m *mockConn) LocalAddr() net.Addr { return nil }
func (m *mockConn) RemoteAddr() net.Addr { return nil }
func (m *mockConn) SetDeadline(t time.Time) error { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestRouterServeTCP(t *testing.T) {
	router := NewRouter()
	router.AddRoute("GET", "/hello", func(req *Request, res *ResponseWriter) {
		res.WriteString(200, "Hello!")
	})

	tests := []struct {
		name           string
		requestPayload []byte
		expectedBody   string
	}{
		{
			name:           "valid request",
			requestPayload: []byte("GET /hello HTTP/1.1\r\nHost: localhost\r\n\r\n"),
			expectedBody:   "Hello!",
		},
		{
			name:           "not found",
			requestPayload: []byte("GET /unknown HTTP/1.1\r\nHost: localhost\r\n\r\n"),
			expectedBody:   "Not Found", 
		},
		{
			name:           "bad request",
			requestPayload: []byte("INVALID_REQUEST_LINE\r\n\r\n"), 
			expectedBody:   "Bad Request", 
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &mockConn{
				readBuffer:  bytes.NewBuffer(tt.requestPayload),
				writeBuffer: &bytes.Buffer{},
			}

			router.ServeTCP(context.Background(), conn)

			response := conn.writeBuffer.String()
			if !bytes.Contains([]byte(response), []byte(tt.expectedBody)) {
				t.Errorf("expected response to contain %q, got: %q", tt.expectedBody, response)
			}
			if !conn.closed {
				t.Errorf("expected connection to be closed")
			}
		})
	}
}
