package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/otel-shop/inventory/internal/store"
)

func TestProductStore_GetStock(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock: %v", err)
		}
		defer func() { _ = db.Close() }()

		mock.ExpectQuery(`SELECT stock FROM products WHERE id = \$1`).
			WithArgs("A123").
			WillReturnRows(sqlmock.NewRows([]string{"stock"}).AddRow(10))

		s := NewProductStore(db)
		got, err := s.GetStock(context.Background(), "A123")
		if err != nil {
			t.Fatalf("GetStock: %v", err)
		}
		if got != 10 {
			t.Fatalf("stock = %d, want 10", got)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock: %v", err)
		}
		defer func() { _ = db.Close() }()

		mock.ExpectQuery(`SELECT stock FROM products WHERE id = \$1`).
			WithArgs("UNKNOWN").
			WillReturnError(sql.ErrNoRows)

		s := NewProductStore(db)
		if _, err := s.GetStock(context.Background(), "UNKNOWN"); !errors.Is(err, store.ErrNotFound) {
			t.Fatalf("err = %v, want store.ErrNotFound", err)
		}
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock: %v", err)
		}
		defer func() { _ = db.Close() }()

		mock.ExpectQuery(`SELECT stock FROM products WHERE id = \$1`).
			WithArgs("A123").
			WillReturnError(errors.New("connection refused"))

		s := NewProductStore(db)
		_, err = s.GetStock(context.Background(), "A123")
		if err == nil || errors.Is(err, store.ErrNotFound) {
			t.Fatalf("err = %v, want generic db error", err)
		}
	})
}

func TestOpen_BadEndpoint(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	database, err := Open(ctx, "postgres://otel:otel@127.0.0.1:1/oteldb?sslmode=disable")
	if err == nil {
		_ = database.Close()
		t.Fatal("Open against a dead endpoint should fail")
	}
}
