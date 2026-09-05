package http

import (
	"net"
	"strconv"
)

var (
	http11     = []byte("HTTP/1.1 ")
	colonSpace = []byte(": ")
)

type ResponseWriter struct {
	conn        net.Conn
	statusCode  int
	headers     []Header
	wroteHeader bool

	//pre-allcated buffer for writing headers to avoid multiple allocations
	headerBuf []byte
}

func (rw *ResponseWriter) Reset(conn net.Conn) {
	rw.conn = conn
	rw.statusCode = 200
	rw.headers = rw.headers[:0]
	rw.wroteHeader = false
	rw.headerBuf = rw.headerBuf[:0]
}

func (rw *ResponseWriter) AddHeader(key, value string) {
	rw.headers = append(rw.headers, Header{
		Key:   []byte(key),
		Value: []byte(value),
	})
}

func (rw *ResponseWriter) SetStatusCode(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) Write(b []byte) (n int, err error) {
	if !rw.wroteHeader {
		rw.AddHeader("Content-Length", strconv.Itoa(len(b)))
		rw.writeHeaders()
	}
	return rw.conn.Write(b)
}

func (rw *ResponseWriter) writeHeaders() {
	if rw.wroteHeader {
		return
	}
	rw.wroteHeader = true

	rw.headerBuf = append(rw.headerBuf, http11...)
	rw.headerBuf = append(rw.headerBuf, strconv.Itoa(rw.statusCode)...)
	rw.headerBuf = append(rw.headerBuf, ' ')
	rw.headerBuf = append(rw.headerBuf, statusText(rw.statusCode)...)
	rw.headerBuf = append(rw.headerBuf, crlf...)

	for _, header := range rw.headers {
		rw.headerBuf = append(rw.headerBuf, header.Key...)
		rw.headerBuf = append(rw.headerBuf, colonSpace...)
		rw.headerBuf = append(rw.headerBuf, header.Value...)
		rw.headerBuf = append(rw.headerBuf, crlf...)
	}

	rw.headerBuf = append(rw.headerBuf, crlf...)

	rw.conn.Write(rw.headerBuf)

}

func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 202:
		return "Accepted"
	case 204:
		return "No Content"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 413:
		return "Payload Too Large"
	case 431:
		return "Request Header Fields Too Large"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}
