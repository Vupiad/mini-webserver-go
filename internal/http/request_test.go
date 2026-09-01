package http

import (
	"bytes"
	"errors"
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name             string
		input            []byte
		expectedError    error
		expectedMethod   []byte
		expectedPath     []byte
		expectedProtocol []byte
		expectedHeaders  []Header
		expectedBody     []byte
	}{
		{
			name:             "Valid GET request",
			input:            []byte("GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"),
			expectedError:    nil,
			expectedMethod:   []byte("GET"),
			expectedPath:     []byte("/"),
			expectedProtocol: []byte("HTTP/1.1"),
			expectedHeaders: []Header{
				{Key: []byte("Host"), Value: []byte("example.com")},
			},
			expectedBody: []byte{},
		},
		{
			name:             "Valid POST request with body",
			input:            []byte("POST /submit HTTP/1.1\r\nHost: example.com\r\nContent-Length: 11\r\n\r\nHello World"),
			expectedError:    nil,
			expectedMethod:   []byte("POST"),
			expectedPath:     []byte("/submit"),
			expectedProtocol: []byte("HTTP/1.1"),
			expectedHeaders: []Header{
				{Key: []byte("Host"), Value: []byte("example.com")},
				{Key: []byte("Content-Length"), Value: []byte("11")},
			},
			expectedBody: []byte("Hello World"),
		},
		{
			name:          "Malformed request line",
			input:         []byte("GET / HTTP/1.1 Host: example.com\r\n\r\n"),
			expectedError: ErrMalformedRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &Request{
				Headers: make([]Header, 0, 20),
			}
			err := ParseRequest(req, tc.input)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}
			if err == nil {
				if !bytes.Equal(req.Method, tc.expectedMethod) {
					t.Errorf("expected method %s, got %s", tc.expectedMethod, req.Method)
				}
				if !bytes.Equal(req.Path, tc.expectedPath) {
					t.Errorf("expected path %s, got %s", tc.expectedPath, req.Path)
				}
				if !bytes.Equal(req.Protocol, tc.expectedProtocol) {
					t.Errorf("expected protocol %s, got %s", tc.expectedProtocol, req.Protocol)
				}
				if len(req.Headers) != len(tc.expectedHeaders) {
					t.Errorf("expected headers length %d, got %d", len(tc.expectedHeaders), len(req.Headers))
				} else {
					for i := range req.Headers {
						if !bytes.Equal(req.Headers[i].Key, tc.expectedHeaders[i].Key) || !bytes.Equal(req.Headers[i].Value, tc.expectedHeaders[i].Value) {
							t.Errorf("expected header %v, got %v", tc.expectedHeaders[i], req.Headers[i])
						}
					}
				}
				if !bytes.Equal(req.ReadBody, tc.expectedBody) {
					t.Errorf("expected body %s, got %s", tc.expectedBody, req.ReadBody)
				}
			}
		})
	}
}

func BenchmarkParseRequest(b *testing.B) {
	rawInput := []byte("GET /users/profile HTTP/1.1\r\nHost: localhost:8080\r\nAccept: application/json\r\n\r\n")

	req := &Request{
		Headers: make([]Header, 0, 20),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ParseRequest(req, rawInput)
	}
}
