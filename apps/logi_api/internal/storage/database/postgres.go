package database

import (
	"context"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB represents a database connection pool.
type DB struct {
	Pool *pgxpool.Pool
}

// NewDatabase creates a new database connection pool.
func NewDatabase(connectionString string) (*DB, error) {
	pool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return &DB{Pool: pool}, nil
}

// Close closes the database connection pool.
func (db *DB) Close() {
	db.Pool.Close()
}

// Executor interface methods - delegate to the pool
func (db *DB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return db.Pool.Query(ctx, sql, args...)
}

func (db *DB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return db.Pool.QueryRow(ctx, sql, args...)
}

func (db *DB) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return db.Pool.Exec(ctx, sql, arguments...)
}

// Transactor interface methods - delegate to the pool
func (db *DB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.Pool.Begin(ctx)
}
