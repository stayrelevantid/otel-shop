package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/otel-shop/inventory/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "18081"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /inventory/{id}", handler.Inventory)

	addr := ":" + port
	log.Printf("inventory-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
