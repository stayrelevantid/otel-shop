package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/otel-shop/inventory/internal/store"
)

// ProductStore implements store.StockStore against PostgreSQL.
type ProductStore struct {
	DB *sql.DB
}

func NewProductStore(database *sql.DB) *ProductStore {
	return &ProductStore{DB: database}
}

// GetStock reads the stock of a product by id (PRD §11).
func (p *ProductStore) GetStock(ctx context.Context, id string) (int, error) {
	var stock int
	err := p.DB.QueryRowContext(ctx,
		`SELECT stock FROM products WHERE id = $1`, id,
	).Scan(&stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, store.ErrNotFound
		}
		return 0, err
	}
	return stock, nil
}
