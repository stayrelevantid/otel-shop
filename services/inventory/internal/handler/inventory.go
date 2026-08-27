package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/otel-shop/inventory/internal/model"
)

// Inventory handles GET /inventory/{id}. Day 2 stub: always returns
// stock 10 — real DB query arrives in Fase 3.
func Inventory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing product id"})
		return
	}

	log.Printf("inventory lookup: id=%s", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.InventoryResponse{Stock: 10})
}
