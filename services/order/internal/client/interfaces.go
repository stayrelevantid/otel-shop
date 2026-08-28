package client

import (
	"context"

	"github.com/otel-shop/order/internal/model"
)

// InventoryClient reads product stock for the checkout flow (F4).
type InventoryClient interface {
	GetStock(ctx context.Context, itemID string) (int, error)
}

// PaymentClient charges a payment for the checkout flow (F4).
type PaymentClient interface {
	Pay(ctx context.Context, req model.PayRequest) error
}
