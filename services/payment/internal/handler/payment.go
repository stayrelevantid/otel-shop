package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/otel-shop/payment/internal/model"
)

// Pay handles POST /pay. Day 2 stub: always success — chaos comes in Fase 5.
func Pay(w http.ResponseWriter, r *http.Request) {
	var req model.PayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OrderID == "" || req.Amount <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.PayResponse{Status: "failed"})
		return
	}

	log.Printf("payment processed: order_id=%s amount=%.2f", req.OrderID, req.Amount)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.PayResponse{Status: "success"})
}
