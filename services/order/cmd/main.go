package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/otel-shop/order/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "18080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /checkout", handler.Checkout)

	addr := ":" + port
	log.Printf("order-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
