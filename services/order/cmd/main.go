package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/otel-shop/order/internal/client"
	"github.com/otel-shop/order/internal/handler"
)

func main() {
	port := envOr("PORT", "18080")
	invURL := envOr("INVENTORY_URL", "http://localhost:18081")
	payURL := envOr("PAYMENT_URL", "http://localhost:18082")

	h := &handler.CheckoutHandler{
		Inv: client.NewInventoryClient(invURL),
		Pay: client.NewPaymentClient(payURL),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /checkout", h.Checkout)

	addr := ":" + port
	log.Printf("order-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
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
