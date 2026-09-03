package db

import (
	"context"
	"database/sql"

	"github.com/XSAM/otelsql"
	_ "github.com/jackc/pgx/v5/stdlib"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Open creates a pooled, instrumented *sql.DB using the pgx stdlib driver.
// Every query emits a DB span (db.system = postgresql) when a TracerProvider
// is installed (F9).
func Open(ctx context.Context, databaseURL string) (*sql.DB, error) {
	database, err := otelsql.Open("pgx", databaseURL,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
	)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(5)

	if err := database.PingContext(ctx); err != nil {
		_ = database.Close()
		return nil, err
	}
	return database, nil
}
