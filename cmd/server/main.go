package main

import (
	"context"
	"fmt"
	"github.com/Vupiad/mini-webserver-go/internal/http"
	"github.com/Vupiad/mini-webserver-go/internal/tcp"
)

func main() {
	router := http.NewRouter()

	router.AddRoute("/ping", func(req *http.Request, res *http.ResponseWriter) {
		res.WriteString(200, "PONG")
	})

	router.AddRoute("/profile", func(req *http.Request, res *http.ResponseWriter) {
		res.WriteString(200, "Profile Data")
	})

	server := &tcp.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Starting server on :8080")
	if err := server.Start(context.Background()); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
