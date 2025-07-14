package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// QueryBuilder helps build common SQL queries
type QueryBuilder struct {
	tableName string
	fields    []string
	values    []interface{}
	where     []string
	orderBy   string
	limit     int
	offset    int
}

// NewQueryBuilder creates a new query builder for the given table
func NewQueryBuilder(tableName string) *QueryBuilder {
	return &QueryBuilder{
		tableName: tableName,
		fields:    make([]string, 0),
		values:    make([]interface{}, 0),
		where:     make([]string, 0),
	}
}

// Select adds fields to select
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.fields = append(qb.fields, fields...)
	return qb
}

// Where adds a WHERE condition
func (qb *QueryBuilder) Where(condition string, value interface{}) *QueryBuilder {
	qb.where = append(qb.where, fmt.Sprintf("%s = $%d", condition, len(qb.values)+1))
	qb.values = append(qb.values, value)
	return qb
}

// OrderBy sets the ORDER BY clause
func (qb *QueryBuilder) OrderBy(field string) *QueryBuilder {
	qb.orderBy = field
	return qb
}

// Limit sets the LIMIT clause
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset sets the OFFSET clause
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// BuildSelect builds a SELECT query
func (qb *QueryBuilder) BuildSelect() (string, []interface{}) {
	fields := "*"
	if len(qb.fields) > 0 {
		fields = strings.Join(qb.fields, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s", fields, qb.tableName)

	if len(qb.where) > 0 {
		query += " WHERE " + strings.Join(qb.where, " AND ")
	}

	if qb.orderBy != "" {
		query += " ORDER BY " + qb.orderBy
	}

	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offset)
	}

	return query, qb.values
}

// DatabaseHelper provides common database operations
type DatabaseHelper struct {
	pool *pgxpool.Pool
}

// NewDatabaseHelper creates a new database helper
func NewDatabaseHelper(pool *pgxpool.Pool) *DatabaseHelper {
	return &DatabaseHelper{pool: pool}
}

// Exists checks if a record exists with the given condition
func (dh *DatabaseHelper) Exists(ctx context.Context, tableName, field string, value interface{}) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1)", tableName, field)
	var exists bool
	err := dh.pool.QueryRow(ctx, query, value).Scan(&exists)
	return exists, err
}

// Count returns the number of records matching the condition
func (dh *DatabaseHelper) Count(ctx context.Context, tableName, field string, value interface{}) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = $1", tableName, field)
	var count int
	err := dh.pool.QueryRow(ctx, query, value).Scan(&count)
	return count, err
}

// WithTransaction executes a function within a database transaction
func (dh *DatabaseHelper) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := dh.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
