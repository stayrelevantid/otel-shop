package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/otel-shop/inventory/internal/store"
)

type fakeStore struct {
	stock int
	err   error
}

func (f fakeStore) GetStock(_ context.Context, _ string) (int, error) {
	return f.stock, f.err
}

func serveInventory(t *testing.T, h *Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /inventory/{id}", h.Get)

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		store      fakeStore
		wantStatus int
		wantStock  int
	}{
		{
			name:       "product exists",
			path:       "/inventory/A123",
			store:      fakeStore{stock: 10},
			wantStatus: http.StatusOK,
			wantStock:  10,
		},
		{
			name:       "product not found",
			path:       "/inventory/UNKNOWN",
			store:      fakeStore{err: store.ErrNotFound},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "database error",
			path:       "/inventory/A123",
			store:      fakeStore{err: errors.New("connection refused")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{Store: tt.store}
			rec := serveInventory(t, h, tt.path)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}

			if tt.wantStatus == http.StatusOK {
				var resp struct {
					Stock int `json:"stock"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Fatalf("decode: %v", err)
				}
				if resp.Stock != tt.wantStock {
					t.Fatalf("stock = %d, want %d", resp.Stock, tt.wantStock)
				}
			}
		})
	}
}
