package http

import (
	"bytes"
	"strings"
)

type node struct {
	children   map[string]*node
	paramChild *node
	paramName  string
	handlers   map[string]Handler
}

func newNode() *node {
	return &node{
		children: make(map[string]*node),
	}
}

func (n *node) insert(method string, path string, handler Handler) {
	curr := n
	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")

	for _, segment := range segments {
		if segment == "" {
			continue
		}

		if strings.HasPrefix(segment, ":") {
			if curr.paramChild == nil {
				curr.paramChild = newNode()
				curr.paramChild.paramName = segment[1:]
			}
			curr = curr.paramChild
		} else {
			if curr.children[segment] == nil {
				curr.children[segment] = newNode()
			}
			curr = curr.children[segment]
		}
	}
	if curr.handlers == nil {
		curr.handlers = make(map[string]Handler)
	}
	curr.handlers[method] = handler
}

func (n *node) search(method []byte, path []byte, req *Request) (Handler, bool) {
	curr := n

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	for len(path) > 0 {
		var segment []byte
		slashIdx := bytes.IndexByte(path, '/')
		if slashIdx == -1 {
			segment = path
			path = nil
		} else {
			segment = path[:slashIdx]
			path = path[slashIdx+1:]
		}

		if len(segment) == 0 {
			continue
		}

		if next, exists := curr.children[string(segment)]; exists {
			curr = next
			continue
		}

		if curr.paramChild != nil {
			curr = curr.paramChild
			if req.Params == nil {
				req.Params = make(map[string][]byte)
			}
			req.Params[curr.paramName] = segment
		}

		return nil, false

	}

	if curr.handlers != nil {
		if handler, exists := curr.handlers[string(method)]; exists {
			return handler, true
		}
	}

	return nil, false
}
