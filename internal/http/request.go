package http

import (
	"bytes"
	"errors"
	"io"
	"net"
)

type Header struct {
	Key   []byte
	Value []byte
}

type Request struct {
	Method   []byte
	Path     []byte
	Protocol []byte
	Headers  []Header

	ContentLength uint32
	ReadBody      []byte
	conn          net.Conn
	BodyBytesRead uint32
}

var (
	ErrMalformedRequest = errors.New("malformed form HTTP request")
	crlf                = []byte("\r\n")
)

func (req *Request) Reset() {
	req.Method = nil
	req.Path = nil
	req.Protocol = nil
	req.Headers = req.Headers[:0]
	req.ContentLength = 0
	req.ReadBody = nil
	req.BodyBytesRead = 0
	req.conn = nil

}

func (req *Request) Read(p []byte) (int, error) {
	if req.BodyBytesRead >= req.ContentLength {
		return 0, io.EOF
	}

	if len(req.ReadBody) > 0 {
		n := copy(p, req.ReadBody)
		req.ReadBody = req.ReadBody[n:]
		req.BodyBytesRead += uint32(n)
		return n, nil
	}

	bytesLeft := req.ContentLength - req.BodyBytesRead

	if len(p) > int(bytesLeft) {
		p = p[:bytesLeft]
	}

	n, err := req.conn.Read(p)
	req.BodyBytesRead += uint32(n)

	return n, err
}
func ParseRequest(req *Request, data []byte) error {
	reqLineEnd := bytes.Index(data, crlf)
	if reqLineEnd == -1 {
		return ErrMalformedRequest
	}

	err := parseRequestLine(req, data[:reqLineEnd])
	if err != nil {
		return err
	}

	remainBytes := data[reqLineEnd+2:]

	req.Headers = req.Headers[:0]
	for {
		if len(remainBytes) >= 2 && remainBytes[0] == '\r' && remainBytes[1] == '\n' {
			req.ReadBody = remainBytes[2:]
			contentLengthErr := extractContentLength(req)
			if contentLengthErr != nil {
				return contentLengthErr
			}
			break
		}
		headerLineEnd := bytes.Index(remainBytes, crlf)
		if headerLineEnd == -1 {
			return ErrMalformedRequest
		}
		err := parseHeaderLine(req, remainBytes[:headerLineEnd])
		if err != nil {
			return err
		}
		remainBytes = remainBytes[headerLineEnd+2:]
	}

	return nil
}

func parseRequestLine(req *Request, line []byte) error {
	space1 := bytes.IndexByte(line, ' ')
	if space1 == -1 {
		return ErrMalformedRequest
	}
	req.Method = line[:space1]

	space2 := bytes.IndexByte(line[space1+1:], ' ')
	if space2 == -1 {
		return ErrMalformedRequest
	}
	req.Path = line[space1+1 : space1+1+space2]
	req.Protocol = line[space1+1+space2+1:]

	if !bytes.Equal(req.Protocol, []byte("HTTP/1.1")) && !bytes.Equal(req.Protocol, []byte("HTTP/1.0")) {
		return ErrMalformedRequest
	}

	return nil
}

func parseHeaderLine(req *Request, line []byte) error {
	colon := bytes.IndexByte(line, ':')
	if colon == -1 {
		return ErrMalformedRequest
	}

	key := bytes.TrimSpace(line[:colon])
	value := bytes.TrimSpace(line[colon+1:])

	req.Headers = append(req.Headers, Header{
		Key:   key,
		Value: value,
	})

	return nil
}

func extractContentLength(req *Request) error {
	for _, header := range req.Headers {
		if bytes.EqualFold(header.Key, []byte("Content-Length")) {
			length, err := parseAsciiInt(header.Value)
			if err != nil {
				return err
			}
			req.ContentLength = length
			return nil
		}
	}
	return nil
}

func parseAsciiInt(b []byte) (uint32, error) {
	var n uint32
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, ErrMalformedRequest
		}
		n = n*10 + uint32(c-'0')
	}
	return n, nil
}
