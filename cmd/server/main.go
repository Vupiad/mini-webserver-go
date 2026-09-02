package main

import (
	"context"
	"fmt"
	"github.com/Vupiad/mini-webserver-go/internal/http"
	"github.com/Vupiad/mini-webserver-go/internal/tcp"
)

func main() {
	router := http.NewRouter()

	router.AddRoute("GET", "/", func(req *http.Request, rw *http.ResponseWriter) {
		rw.WriteString(200, "Hello, World!")
	})
	router.AddRoute("GET", "/hello/:name", func(req *http.Request, rw *http.ResponseWriter) {
		name := req.Params["name"]
		rw.WriteString(200, fmt.Sprintf("Hello, %s!", name))
	})

	Server := &tcp.Server{
		Addr:    ":8080",
		Handler: router,
	}

	ctx := context.Background()
	if err := Server.Start(ctx); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

}
