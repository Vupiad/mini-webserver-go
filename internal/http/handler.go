package http

import (
	"fmt"
	"net"
)

type ResponseWriter struct {
	Conn net.Conn
}

// just mock the response writer for now, we will implement it later
func (rw *ResponseWriter) WriteString(statusCode int, body string) {
	response := fmt.Sprintf("HTTP/1.1 %d OK\r\n"+
		"Content-Length: %d\r\n"+
		"Content-Type: text/plain\r\n"+
		"Connection: close\r\n"+
		"\r\n"+
		"%s", statusCode, len(body), body)

	rw.Conn.Write([]byte(response))
}

type Handler interface {
	ServeHTTP(req *Request, rw *ResponseWriter)
}

type HandlerFunc func(req *Request, rw *ResponseWriter)

func (f HandlerFunc) ServeHTTP(req *Request, rw *ResponseWriter) {
	f(req, rw)
}
