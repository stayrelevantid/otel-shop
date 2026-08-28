package store

import (
	"context"
	"errors"
)

// ErrNotFound is returned when a product does not exist.
var ErrNotFound = errors.New("product not found")

// StockStore abstracts product stock lookup so handlers stay testable.
type StockStore interface {
	GetStock(ctx context.Context, id string) (int, error)
}
