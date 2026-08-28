package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/otel-shop/order/internal/client"
	"github.com/otel-shop/order/internal/model"
)

// CheckoutHandler serves POST /checkout with injected downstream clients (F4).
type CheckoutHandler struct {
	Inv client.InventoryClient
	Pay client.PaymentClient
}

// Checkout validates the request, checks stock, charges payment, and returns.
func (h *CheckoutHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req model.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ItemID == "" || req.Qty <= 0 {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "item_id is required and qty must be positive"})
		return
	}

	ctx := r.Context()

	stock, err := h.Inv.GetStock(ctx, req.ItemID)
	if err != nil {
		log.Printf("checkout: inventory error for %s: %v", req.ItemID, err)
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "inventory failed"})
		return
	}
	if stock < req.Qty {
		log.Printf("checkout: insufficient stock %s have=%d want=%d", req.ItemID, stock, req.Qty)
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "insufficient stock"})
		return
	}

	orderID := newOrderID()
	if err := h.Pay.Pay(ctx, model.PayRequest{OrderID: orderID, Amount: float64(req.Qty) * 10}); err != nil {
		log.Printf("checkout: payment error for %s: %v", orderID, err)
		respondJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "payment failed"})
		return
	}

	log.Printf("checkout: ok order_id=%s item=%s qty=%d", orderID, req.ItemID, req.Qty)
	respondJSON(w, http.StatusOK, model.CheckoutResponse{OrderID: orderID, Status: "paid"})
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
