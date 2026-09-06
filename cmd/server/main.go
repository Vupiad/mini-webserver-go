package main

import (
	"context"
	"fmt"
	"github.com/Vupiad/mini-webserver-go/internal/http"
	"github.com/Vupiad/mini-webserver-go/internal/tcp"
)

func main() {
	router := http.NewRouter()

	router.Use(http.Logger())
	router.AddRoute("GET", "/", func(req *http.Request, rw *http.ResponseWriter) {
		rw.Write([]byte("Hello, World!"))
	})

	router.AddRoute("GET", "/servername", func(req *http.Request, rw *http.ResponseWriter) {
		rw.Write([]byte("Server Name: My Go Server"))
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
