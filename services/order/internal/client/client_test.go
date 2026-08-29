package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/otel-shop/order/internal/model"
)

func TestInventoryClient_GetStock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/inventory/A123" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"stock":10}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := NewInventoryClient(srv.URL)

	stock, err := c.GetStock(context.Background(), "A123")
	if err != nil {
		t.Fatalf("GetStock(A123): %v", err)
	}
	if stock != 10 {
		t.Fatalf("stock = %d, want 10", stock)
	}

	if _, err := c.GetStock(context.Background(), "UNKNOWN"); err == nil {
		t.Fatal("GetStock(UNKNOWN) should fail")
	}
}

func TestPaymentClient_Pay(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var req model.PayRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if req.OrderID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewPaymentClient(srv.URL)

	if err := c.Pay(context.Background(), model.PayRequest{OrderID: "O-1", Amount: 10}); err != nil {
		t.Fatalf("Pay(valid): %v", err)
	}
	if err := c.Pay(context.Background(), model.PayRequest{OrderID: "", Amount: 10}); err == nil {
		t.Fatal("Pay(invalid) should fail")
	}
}

func TestPaymentClient_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewPaymentClient(srv.URL)
	if err := c.Pay(context.Background(), model.PayRequest{OrderID: "O-1", Amount: 10}); err == nil {
		t.Fatal("Pay against 500 should fail")
	}
}
