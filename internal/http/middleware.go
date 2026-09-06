package http

import (
	"fmt"
	"time"
)

func Logger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(req *Request, rw *ResponseWriter) {
			start := time.Now()
			next(req, rw)
			duration := time.Since(start)

			fmt.Printf("[LOG] %s %s | Status: %d | Time: %v\n",
				string(req.Method),
				string(req.Path),
				rw.statusCode,
				duration,
			)
		}
	}
}
