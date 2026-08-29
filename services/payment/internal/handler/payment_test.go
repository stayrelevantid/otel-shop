package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/otel-shop/payment/internal/chaos"
)

func servePay(t *testing.T, h *Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/pay", bytes.NewReader([]byte(body)))
	rec := httptest.NewRecorder()
	h.Pay(rec, req)
	return rec
}

func TestPay(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		chaos    chaos.Config
		wantCode int
		wantBody string
	}{
		{
			name:     "successful payment",
			body:     `{"order_id":"O-1","amount":100}`,
			chaos:    chaos.Config{},
			wantCode: http.StatusOK,
			wantBody: "success",
		},
		{
			name:     "forced chaos error",
			body:     `{"order_id":"O-1","amount":100}`,
			chaos:    chaos.Config{ErrorPercent: 100},
			wantCode: http.StatusInternalServerError,
			wantBody: "failed",
		},
		{
			name:     "invalid json",
			body:     `{`,
			chaos:    chaos.Config{},
			wantCode: http.StatusBadRequest,
			wantBody: "failed",
		},
		{
			name:     "negative amount",
			body:     `{"order_id":"O-1","amount":-5}`,
			chaos:    chaos.Config{},
			wantCode: http.StatusBadRequest,
			wantBody: "failed",
		},
		{
			name:     "missing order id",
			body:     `{"order_id":"","amount":100}`,
			chaos:    chaos.Config{},
			wantCode: http.StatusBadRequest,
			wantBody: "failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{Chaos: tt.chaos}
			rec := servePay(t, h, tt.body)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantCode, rec.Body.String())
			}
			if !bytes.Contains(rec.Body.Bytes(), []byte(tt.wantBody)) {
				t.Fatalf("body = %s, want containing %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}
