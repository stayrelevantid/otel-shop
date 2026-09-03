package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"

	"github.com/otel-shop/inventory/internal/model"
	"github.com/otel-shop/inventory/internal/store"
)

// Handler serves the inventory API with an injected StockStore.
type Handler struct {
	Store store.StockStore
}

// Get handles GET /inventory/{id} (PRD §7). Span attributes (F11.2) and
// baggage readout (F11.6) are attached to the otelhttp server span.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondNotFound(w, "missing product id")
		return
	}

	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.String("product.id", id))
	if m := baggage.FromContext(r.Context()).Member("order.id"); m.Key() != "" {
		span.SetAttributes(attribute.String("baggage.order.id", m.Value()))
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

	span.SetAttributes(attribute.Int("product.stock", stock))
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
