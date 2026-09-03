package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/otel-shop/payment/internal/chaos"
	"github.com/otel-shop/payment/internal/model"
)

// Handler serves the payment API with configurable chaos.
type Handler struct {
	Chaos chaos.Config
}

// Pay handles POST /pay (PRD §8). Chaos is applied before a successful charge.
// Span attributes (F11.3), events (F11.4), and error status (F11.5) are
// attached to the otelhttp server span.
func (h *Handler) Pay(w http.ResponseWriter, r *http.Request) {
	var req model.PayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OrderID == "" || req.Amount <= 0 {
		respondJSON(w, http.StatusBadRequest, model.PayResponse{Status: "failed"})
		return
	}

	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(
		attribute.String("payment.order_id", req.OrderID),
		attribute.Float64("payment.amount", req.Amount),
	)
	if m := baggage.FromContext(r.Context()).Member("order.id"); m.Key() != "" {
		span.SetAttributes(attribute.String("baggage.order.id", m.Value()))
	}

	span.AddEvent("payment_started")

	if err := h.Chaos.Apply(r.Context()); err != nil {
		log.Printf("payment chaos: %v", err)
		span.AddEvent("payment_failed")
		span.RecordError(err)
		span.SetAttributes(attribute.String("payment.status", "failed"))
		span.SetStatus(codes.Error, "chaos induced payment failure")
		respondJSON(w, http.StatusInternalServerError, model.PayResponse{Status: "failed"})
		return
	}

	log.Printf("payment processed: order_id=%s amount=%.2f", req.OrderID, req.Amount)
	span.AddEvent("payment_completed")
	span.SetAttributes(attribute.String("payment.status", "success"))
	respondJSON(w, http.StatusOK, model.PayResponse{Status: "success"})
}

func respondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
