package http

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type Router struct {
	routes      map[string]Handler
	requestPool sync.Pool
	bufferPool  sync.Pool
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Handler),
		requestPool: sync.Pool{
			New: func() interface{} {
				return &Request{Headers: make([]Header, 0, 20)}
			},
		},
		bufferPool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, 4096)
				return &b
			},
		},
	}
}

func (r *Router) AddRoute(path string, handler HandlerFunc) {
	r.routes[path] = handler
}

func (r *Router) ServeTCP(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	bufferPtr := r.bufferPool.Get().(*[]byte)
	buffer := *bufferPtr
	defer r.bufferPool.Put(bufferPtr)

	req := r.requestPool.Get().(*Request)

	req.conn = conn

	defer func() {
		req.Reset()
		r.requestPool.Put(req)
	}()

	n, err := conn.Read(buffer)

	if err != nil {
		return
	}

	if err := ParseRequest(req, buffer[:n]); err != nil {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 11\r\n\r\nBad Request"))
		return
	}

	res := &ResponseWriter{Conn: conn}

	if handler, ok := r.routes[string(req.Path)]; ok {
		fmt.Println("Handling request for path:", string(req.Path))
		handler.ServeHTTP(req, res)
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 9\r\n\r\nNot Found"))
	}
}
