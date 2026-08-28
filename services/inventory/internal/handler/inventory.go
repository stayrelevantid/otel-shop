package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/otel-shop/inventory/internal/model"
	"github.com/otel-shop/inventory/internal/store"
)

// Handler serves the inventory API with an injected StockStore.
type Handler struct {
	Store store.StockStore
}

// Get handles GET /inventory/{id} (PRD §7).
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondNotFound(w, "missing product id")
		return
	}

	stock, err := h.Store.GetStock(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			log.Printf("inventory: product %s not found", id)
			respondNotFound(w, "product not found")
			return
		}
		log.Printf("inventory: db error for %s: %v", id, err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	log.Printf("inventory lookup: id=%s stock=%d", id, stock)
	respondJSON(w, http.StatusOK, model.InventoryResponse{Stock: stock})
}

func respondNotFound(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func respondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
