package http

import (
	"testing"
)

type mockHandler struct {
	name string
}

func (m *mockHandler) ServeHTTP(req *Request, res *ResponseWriter) {}

func TestTrie(t *testing.T) {
	root := newNode()
	handlerA := &mockHandler{"A"}
	handlerB := &mockHandler{"B"}
	handlerC := &mockHandler{"C"}

	root.insert("GET", "/", handlerA)
	root.insert("GET", "/about", handlerB)
	root.insert("GET", "/users/:id", handlerC)

	tests := []struct {
		name        string
		method      string
		path        string
		wantHandler Handler
		wantMatch   bool
		wantParams  map[string]string
	}{
		{
			name:        "root path",
			method:      "GET",
			path:        "/",
			wantHandler: handlerA,
			wantMatch:   true,
		},
		{
			name:        "static path",
			method:      "GET",
			path:        "/about",
			wantHandler: handlerB,
			wantMatch:   true,
		},
		{
			name:        "param path",
			method:      "GET",
			path:        "/users/123",
			wantHandler: handlerC,
			wantMatch:   true,
			wantParams:  map[string]string{"id": "123"},
		},
		{
			name:        "not found",
			method:      "GET",
			path:        "/unknown",
			wantHandler: nil,
			wantMatch:   false,
		},
		{
			name:        "method not allowed",
			method:      "POST",
			path:        "/",
			wantHandler: nil,
			wantMatch:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{}
			handler, match := root.search([]byte(tt.method), []byte(tt.path), req)

			if match != tt.wantMatch {
				t.Errorf("expected match %v, got %v", tt.wantMatch, match)
			}
			if handler != tt.wantHandler {
				t.Errorf("expected handler %v, got %v", tt.wantHandler, handler)
			}
			
			if tt.wantParams != nil {
				if len(req.Params) != len(tt.wantParams) {
					t.Errorf("expected %d params, got %d", len(tt.wantParams), len(req.Params))
				}
				for k, v := range tt.wantParams {
					if string(req.Params[k]) != v {
						t.Errorf("expected param %s to be %s, got %s", k, v, string(req.Params[k]))
					}
				}
			}
		})
	}
}
