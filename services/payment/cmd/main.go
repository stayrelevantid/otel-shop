package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/otel-shop/payment/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "18082"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /pay", handler.Pay)

	addr := ":" + port
	log.Printf("payment-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
