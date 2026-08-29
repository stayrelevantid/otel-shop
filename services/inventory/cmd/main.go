package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	inventorydb "github.com/otel-shop/inventory/internal/db"
	"github.com/otel-shop/inventory/internal/handler"
	"github.com/otel-shop/telemetry"
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

	tctx, tcancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer tcancel()

	shutdown, err := telemetry.Init(tctx, "inventory-service", "1.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "telemetry init: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("telemetry shutdown: %v", err)
		}
	}()

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
	if err := http.ListenAndServe(addr, otelhttp.NewHandler(mux, "inventory-service")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
