package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	inventorydb "github.com/otel-shop/inventory/internal/db"
	"github.com/otel-shop/inventory/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "18081"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is required")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	database, err := inventorydb.Open(ctx, databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = database.Close() }()

	h := &handler.Handler{
		Store: inventorydb.NewProductStore(database),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /inventory/{id}", h.Get)

	addr := ":" + port
	log.Printf("inventory-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
