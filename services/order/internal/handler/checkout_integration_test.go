package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/otel-shop/order/internal/client"
)

// TestCheckoutIntegration wires the real CheckoutHandler with the real HTTP
// clients against httptest servers acting as Inventory and Payment —
// a full-stack test minus the cluster (F13.1).
func TestCheckoutIntegration(t *testing.T) {
	invSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/inventory/A123" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"stock":10}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer invSrv.Close()

	paySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			OrderID string  `json:"order_id"`
			Amount  float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OrderID == "" || req.Amount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success"}`))
	}))
	defer paySrv.Close()

	h := &CheckoutHandler{
		Inv: client.NewInventoryClient(invSrv.URL),
		Pay: client.NewPaymentClient(paySrv.URL),
	}

	post := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewReader([]byte(body)))
		rec := httptest.NewRecorder()
		h.Checkout(rec, req)
		return rec
	}

	t.Run("full flow success", func(t *testing.T) {
		rec := post(`{"item_id":"A123","qty":2}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
		}
		var resp struct {
			OrderID string `json:"order_id"`
			Status  string `json:"status"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if resp.Status != "paid" || resp.OrderID == "" {
			t.Fatalf("resp = %+v, want paid with order id", resp)
		}
	})

	t.Run("inventory not found propagates as 500", func(t *testing.T) {
		rec := post(`{"item_id":"UNKNOWN","qty":1}`)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500 (body: %s)", rec.Code, rec.Body.String())
		}
	})

	t.Run("insufficient stock propagates as 400", func(t *testing.T) {
		rec := post(`{"item_id":"A123","qty":99}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400 (body: %s)", rec.Code, rec.Body.String())
		}
	})
}
