package http

type Handler interface {
	ServeHTTP(req *Request, rw *ResponseWriter)
}

type HandlerFunc func(req *Request, rw *ResponseWriter)

type Middleware func(HandlerFunc) HandlerFunc

func (f HandlerFunc) ServeHTTP(req *Request, rw *ResponseWriter) {
	f(req, rw)
}
