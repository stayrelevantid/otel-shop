package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/otel-shop/order/internal/client"
	"github.com/otel-shop/order/internal/handler"
	"github.com/otel-shop/telemetry"
)

func main() {
	port := envOr("PORT", "18080")
	invURL := envOr("INVENTORY_URL", "http://localhost:18081")
	payURL := envOr("PAYMENT_URL", "http://localhost:18082")

	ctx := context.Background()
	shutdown, err := telemetry.Init(ctx, "order-service", "1.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "telemetry init: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("telemetry shutdown: %v", err)
		}
	}()

	h := &handler.CheckoutHandler{
		Inv: client.NewInventoryClient(invURL),
		Pay: client.NewPaymentClient(payURL),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /checkout", h.Checkout)

	addr := ":" + port
	log.Printf("order-service listening on %s", addr)
	if err := http.ListenAndServe(addr, otelhttp.NewHandler(mux, "order-service")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
