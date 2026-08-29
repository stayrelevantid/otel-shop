package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/otel-shop/order/internal/model"
)

type fakeInventory struct {
	stock int
	err   error
}

func (f fakeInventory) GetStock(_ context.Context, _ string) (int, error) {
	return f.stock, f.err
}

type fakePayment struct {
	err error
}

func (f fakePayment) Pay(_ context.Context, _ model.PayRequest) error {
	return f.err
}

func TestCheckout(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		inv        fakeInventory
		pay        fakePayment
		wantStatus int
		wantErr    string
	}{
		{
			name:       "success",
			body:       `{"item_id":"A123","qty":1}`,
			inv:        fakeInventory{stock: 10},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid json",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "item_id is required and qty must be positive",
		},
		{
			name:       "invalid qty",
			body:       `{"item_id":"A123","qty":0}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    "item_id is required and qty must be positive",
		},
		{
			name:       "inventory failure",
			body:       `{"item_id":"A123","qty":1}`,
			inv:        fakeInventory{err: errors.New("db down")},
			wantStatus: http.StatusInternalServerError,
			wantErr:    "inventory failed",
		},
		{
			name:       "insufficient stock",
			body:       `{"item_id":"C123","qty":999}`,
			inv:        fakeInventory{stock: 5},
			wantStatus: http.StatusBadRequest,
			wantErr:    "insufficient stock",
		},
		{
			name:       "payment failure",
			body:       `{"item_id":"A123","qty":1}`,
			inv:        fakeInventory{stock: 10},
			pay:        fakePayment{err: errors.New("chaos")},
			wantStatus: http.StatusInternalServerError,
			wantErr:    "payment failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CheckoutHandler{Inv: tt.inv, Pay: tt.pay}
			req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewReader([]byte(tt.body)))
			rec := httptest.NewRecorder()

			h.Checkout(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}

			var resp struct {
				OrderID string `json:"order_id"`
				Status  string `json:"status"`
				Error   string `json:"error"`
			}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			if tt.wantErr != "" {
				if resp.Error != tt.wantErr {
					t.Fatalf("error = %q, want %q", resp.Error, tt.wantErr)
				}
				return
			}
			if resp.Status != "paid" {
				t.Fatalf("status = %q, want %q", resp.Status, "paid")
			}
			if resp.OrderID == "" {
				t.Fatal("order_id is empty")
			}
		})
	}
}

func TestNewOrderID(t *testing.T) {
	a, b := newOrderID(), newOrderID()
	if a == b {
		t.Fatalf("order ids should differ, got %q twice", a)
	}
	if len(a) < 3 || a[:2] != "O-" {
		t.Fatalf("order id %q should have O- prefix", a)
	}
}
