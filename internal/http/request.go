package http

import (
	"bytes"
	"errors"
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
	Body          []byte
}

var (
	ErrMalformedRequest = errors.New("malformed form HTTP request")
	crlf                = []byte("\r\n")
)

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

	for {
		if remainBytes[0] == '\r' && remainBytes[1] == '\n' {
			req.Body = remainBytes[2:]
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
