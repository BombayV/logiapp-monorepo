package database

import (
	"context"
	"log"

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
