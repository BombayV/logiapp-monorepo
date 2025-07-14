package database

import (
	"bombayv/logiapp-monorepo/logi_api/internal/core/orders"
	"context"
	"fmt"
)

// OrderRepository handles the database operations for orders.
type OrderRepository struct {
	db     *DB
	helper *DatabaseHelper
}

// NewOrderRepository creates a new OrderRepository.
func NewOrderRepository(db *DB) *OrderRepository {
	return &OrderRepository{
		db:     db,
		helper: NewDatabaseHelper(db.Pool),
	}
}

// Save saves an order to the database.
func (r *OrderRepository) Save(ctx context.Context, o *orders.Order) error {
	query := `
		INSERT INTO orders (order_id, created_by, assigned_to, delivery_address, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Pool.Exec(ctx, query, o.OrderID, o.CreatedBy, o.AssignedTo, o.DeliveryAddress, o.Status, o.CreatedAt, o.UpdatedAt)
	return err
}

// FindAll returns all orders
func (r *OrderRepository) FindAll(ctx context.Context) ([]*orders.Order, error) {
	return r.findOrdersWithQuery(ctx, "SELECT order_id, created_by, assigned_to, delivery_address, status, created_at, updated_at FROM orders ORDER BY created_at DESC")
}

// FindByID finds an order by its ID
func (r *OrderRepository) FindByID(ctx context.Context, orderID string) (*orders.Order, error) {
	query := `
		SELECT order_id, created_by, assigned_to, delivery_address, status, created_at, updated_at
		FROM orders
		WHERE order_id = $1
	`
	var o orders.Order
	err := r.db.Pool.QueryRow(ctx, query, orderID).Scan(
		&o.OrderID, &o.CreatedBy, &o.AssignedTo, &o.DeliveryAddress,
		&o.Status, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// FindByUserID finds all orders for a specific user
func (r *OrderRepository) FindByUserID(ctx context.Context, userID string) ([]*orders.Order, error) {
	query := `
		SELECT order_id, created_by, assigned_to, delivery_address, status, created_at, updated_at
		FROM orders
		WHERE created_by = $1
		ORDER BY created_at DESC
	`
	return r.findOrdersWithQuery(ctx, query, userID)
}

// FindByStatus finds all orders with a specific status
func (r *OrderRepository) FindByStatus(ctx context.Context, status string) ([]*orders.Order, error) {
	query := `
		SELECT order_id, created_by, assigned_to, delivery_address, status, created_at, updated_at
		FROM orders
		WHERE status = $1
		ORDER BY created_at DESC
	`
	return r.findOrdersWithQuery(ctx, query, status)
}

// FindWithPagination returns orders with pagination
func (r *OrderRepository) FindWithPagination(ctx context.Context, limit, offset int) ([]*orders.Order, int, error) {
	// Get total count
	totalCount, err := r.helper.Count(ctx, "orders", "1", 1) // Count all records
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT order_id, created_by, assigned_to, delivery_address, status, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	ordersList, err := r.findOrdersWithQuery(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return ordersList, totalCount, nil
}

// Update updates an existing order
func (r *OrderRepository) Update(ctx context.Context, o *orders.Order) error {
	query := `
		UPDATE orders 
		SET assigned_to = $2, delivery_address = $3, status = $4, updated_at = $5
		WHERE order_id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, o.OrderID, o.AssignedTo, o.DeliveryAddress, o.Status, o.UpdatedAt)
	return err
}

// Delete deletes an order (hard delete)
func (r *OrderRepository) Delete(ctx context.Context, orderID string) error {
	query := `DELETE FROM orders WHERE order_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, orderID)
	return err
}

// Exists checks if an order exists
func (r *OrderRepository) Exists(ctx context.Context, orderID string) (bool, error) {
	return r.helper.Exists(ctx, "orders", "order_id", orderID)
}

// CountByStatus returns the count of orders by status
func (r *OrderRepository) CountByStatus(ctx context.Context, status string) (int, error) {
	return r.helper.Count(ctx, "orders", "status", status)
}

// findOrdersWithQuery is a helper method to execute queries that return multiple orders
func (r *OrderRepository) findOrdersWithQuery(ctx context.Context, query string, args ...interface{}) ([]*orders.Order, error) {
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordersList []*orders.Order
	for rows.Next() {
		var o orders.Order
		if err := rows.Scan(&o.OrderID, &o.CreatedBy, &o.AssignedTo, &o.DeliveryAddress, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		ordersList = append(ordersList, &o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return ordersList, nil
}
