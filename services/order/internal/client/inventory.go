package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/otel-shop/order/internal/model"
)

// InventoryHTTP is an HTTP implementation of InventoryClient.
type InventoryHTTP struct {
	BaseURL string
	HTTP    *http.Client
}

// NewInventoryClient builds a client for the Inventory service. The transport
// propagates trace context (traceparent + baggage) on every call (F8).
func NewInventoryClient(baseURL string) *InventoryHTTP {
	return &InventoryHTTP{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout:   5 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

// GetStock calls GET /inventory/{id} and returns the stock level.
func (c *InventoryHTTP) GetStock(ctx context.Context, itemID string) (int, error) {
	endpoint := fmt.Sprintf("%s/inventory/%s", c.BaseURL, url.PathEscape(itemID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return 0, fmt.Errorf("inventory not found: %w", model.ErrInventory)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("inventory status %d: %w", resp.StatusCode, model.ErrInventory)
	}

	var r model.InventoryResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return 0, fmt.Errorf("inventory decode: %w", model.ErrInventory)
	}
	return r.Stock, nil
}
