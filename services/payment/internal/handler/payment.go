package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/otel-shop/payment/internal/chaos"
	"github.com/otel-shop/payment/internal/model"
)

// Handler serves the payment API with configurable chaos.
type Handler struct {
	Chaos chaos.Config
}

// Pay handles POST /pay (PRD §8). Chaos is applied before a successful charge.
func (h *Handler) Pay(w http.ResponseWriter, r *http.Request) {
	var req model.PayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OrderID == "" || req.Amount <= 0 {
		respondJSON(w, http.StatusBadRequest, model.PayResponse{Status: "failed"})
		return
	}

	if err := h.Chaos.Apply(r.Context()); err != nil {
		log.Printf("payment chaos: %v", err)
		respondJSON(w, http.StatusInternalServerError, model.PayResponse{Status: "failed"})
		return
	}

	log.Printf("payment processed: order_id=%s amount=%.2f", req.OrderID, req.Amount)
	respondJSON(w, http.StatusOK, model.PayResponse{Status: "success"})
}

func respondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
