package http

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type Router struct {
	root         *node
	requestPool  sync.Pool
	bufferPool   sync.Pool
	responsePool sync.Pool
}

func NewRouter() *Router {
	return &Router{
		root: newNode(),
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
		responsePool: sync.Pool{
			New: func() interface{} {
				return &ResponseWriter{
					headers:   make([]Header, 0, 10),
					headerBuf: make([]byte, 0, 1024)}
			},
		},
	}
}

func (r *Router) AddRoute(method string, path string, handler HandlerFunc) {
	r.root.insert(method, path, handler)
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

	res := r.responsePool.Get().(*ResponseWriter)
	res.Reset(conn)

	defer func() {
		if !res.wroteHeader {
			res.writeHeaders()
		}
		r.responsePool.Put(res)
	}()

	if err := ParseRequest(req, buffer[:n]); err != nil {
		res.SetStatusCode(400)
		res.Write([]byte("Bad Request"))
		return
	}

	if handler, ok := r.root.search(req.Method, req.Path, req); ok {
		fmt.Println("Handling request for path:", string(req.Path))
		handler.ServeHTTP(req, res)
	} else {
		res.SetStatusCode(404)
		res.Write([]byte("Not Found"))
	}
}
