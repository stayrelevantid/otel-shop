package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/otel-shop/order/internal/client"
	"github.com/otel-shop/order/internal/model"
)

var tracer = otel.Tracer("github.com/otel-shop/order")

// CheckoutHandler serves POST /checkout with injected downstream clients (F4).
type CheckoutHandler struct {
	Inv client.InventoryClient
	Pay client.PaymentClient
}

// Checkout validates the request, checks stock, charges payment, and returns.
// Manual spans (F10), baggage (F11.6), and error status (F11.7) are applied
// around each step.
func (h *CheckoutHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req model.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ItemID == "" || req.Qty <= 0 {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "item_id is required and qty must be positive"})
		return
	}

	// The otelhttp server span — used to surface failures on the parent too.
	serverSpan := trace.SpanFromContext(r.Context())

	ctx, span := tracer.Start(r.Context(), "validate-order")
	span.SetAttributes(
		attribute.String("order.item_id", req.ItemID),
		attribute.Int("order.quantity", req.Qty),
	)
	span.End()

	// Generate the order id early so it can ride the baggage to every
	// downstream service (PRD §22).
	orderID := newOrderID()
	if bag, err := baggage.New(); err == nil {
		if member, err := baggage.NewMember("order.id", orderID); err == nil {
			if withID, err := bag.SetMember(member); err == nil {
				ctx = baggage.ContextWithBaggage(ctx, withID)
			}
		}
	}

	// F10.2: check-inventory span.
	ictx, ispan := tracer.Start(ctx, "check-inventory")
	ispan.SetAttributes(attribute.String("product.id", req.ItemID))

	stock, err := h.Inv.GetStock(ictx, req.ItemID)
	if err != nil {
		ispan.RecordError(err)
		ispan.SetStatus(codes.Error, "inventory failed")
		ispan.End()
		markFailure(serverSpan, err, "inventory failed")
		log.Printf("checkout: inventory error for %s: %v", req.ItemID, err)
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "inventory failed"})
		return
	}
	ispan.SetAttributes(attribute.Int("product.stock", stock))
	if stock < req.Qty {
		ispan.SetStatus(codes.Error, "insufficient stock")
		ispan.End()
		markFailure(serverSpan, model.ErrInsufficientStock, "insufficient stock")
		log.Printf("checkout: insufficient stock %s have=%d want=%d", req.ItemID, stock, req.Qty)
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "insufficient stock"})
		return
	}
	ispan.End()

	// F10.3: process-payment span.
	amount := float64(req.Qty) * 10
	pctx, pspan := tracer.Start(ctx, "process-payment")
	pspan.SetAttributes(
		attribute.String("payment.order_id", orderID),
		attribute.Float64("payment.amount", amount),
	)

	if err := h.Pay.Pay(pctx, model.PayRequest{OrderID: orderID, Amount: amount}); err != nil {
		pspan.RecordError(err)
		pspan.SetAttributes(attribute.String("payment.status", "failed"))
		pspan.SetStatus(codes.Error, "payment failed")
		pspan.End()
		markFailure(serverSpan, err, "payment failed")
		log.Printf("checkout: payment error for %s: %v", orderID, err)
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "payment failed"})
		return
	}
	pspan.SetAttributes(attribute.String("payment.status", "success"))
	pspan.End()

	log.Printf("checkout: ok order_id=%s item=%s qty=%d", orderID, req.ItemID, req.Qty)
	respondJSON(w, http.StatusOK, model.CheckoutResponse{OrderID: orderID, Status: "paid"})
}

// markFailure surfaces a downstream failure on the parent server span (F11.7).
func markFailure(span trace.Span, err error, msg string) {
	span.RecordError(err)
	span.SetStatus(codes.Error, msg)
}

func newOrderID() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("O-%d", time.Now().Unix())
	}
	return "O-" + hex.EncodeToString(b)
}

func respondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
