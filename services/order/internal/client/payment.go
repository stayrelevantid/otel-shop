package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/otel-shop/order/internal/model"
)

// PaymentHTTP is an HTTP implementation of PaymentClient.
type PaymentHTTP struct {
	BaseURL string
	HTTP    *http.Client
}

// NewPaymentClient builds a client for the Payment service. The transport
// propagates trace context (traceparent + baggage) on every call (F8).
func NewPaymentClient(baseURL string) *PaymentHTTP {
	return &PaymentHTTP{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout:   5 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

// Pay calls POST /pay and returns an error on any non-200 response.
func (c *PaymentHTTP) Pay(ctx context.Context, req model.PayRequest) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("payment encode: %w", model.ErrPayment)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/pay", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("payment request: %w", model.ErrPayment)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return fmt.Errorf("payment call: %w", model.ErrPayment)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("payment status %d: %s: %w", resp.StatusCode, string(body), model.ErrPayment)
	}
	return nil
}
