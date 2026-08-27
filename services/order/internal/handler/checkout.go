package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/otel-shop/order/internal/model"
)

// Checkout handles POST /checkout. Day 2 stub: validate, generate
// an order id, and respond as paid — no downstream calls yet.
func Checkout(w http.ResponseWriter, r *http.Request) {
	var req model.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	if req.ItemID == "" || req.Qty <= 0 {
		respondJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "item_id is required and qty must be positive"})
		return
	}

	orderID := newOrderID()
	log.Printf("checkout: item=%s qty=%d order_id=%s", req.ItemID, req.Qty, orderID)

	respondJSON(w, http.StatusOK, model.CheckoutResponse{
		OrderID: orderID,
		Status:  "paid",
	})
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
