package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/otel-shop/payment/internal/chaos"
	"github.com/otel-shop/payment/internal/handler"
	"github.com/otel-shop/telemetry"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "18082"
	}

	ctx := context.Background()
	shutdown, err := telemetry.Init(ctx, "payment-service", "1.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "telemetry init: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("telemetry shutdown: %v", err)
		}
	}()

	h := &handler.Handler{
		Chaos: chaos.FromEnv(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /pay", h.Pay)

	addr := ":" + port
	log.Printf("payment-service listening on %s", addr)
	if err := http.ListenAndServe(addr, otelhttp.NewHandler(mux, "payment-service")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
