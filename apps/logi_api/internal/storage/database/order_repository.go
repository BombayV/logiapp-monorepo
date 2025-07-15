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
		INSERT INTO orders (order_id, order_number, created_by, assigned_to, delivery_address, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Pool.Exec(ctx, query, o.OrderID, o.OrderNumber, o.CreatedBy, o.AssignedTo, o.DeliveryAddress, o.Status, o.CreatedAt, o.UpdatedAt)
	return err
}

// FindAll returns all orders
func (r *OrderRepository) FindAll(ctx context.Context) ([]*orders.Order, error) {
	return r.findOrdersWithQuery(ctx, "SELECT order_id, order_number, created_by, assigned_to, delivery_address, status, created_at, updated_at FROM orders ORDER BY created_at DESC")
}

// FindByID finds an order by its ID
func (r *OrderRepository) FindByID(ctx context.Context, orderID string) (*orders.Order, error) {
	query := `
		SELECT order_id, order_number, created_by, assigned_to, delivery_address, status, created_at, updated_at
		FROM orders
		WHERE order_id = $1
	`
	var o orders.Order
	err := r.db.Pool.QueryRow(ctx, query, orderID).Scan(
		&o.OrderID, &o.OrderNumber, &o.CreatedBy, &o.AssignedTo, &o.DeliveryAddress,
		&o.Status, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// FindByIDWithItems finds an order by its ID including all order items
func (r *OrderRepository) FindByIDWithItems(ctx context.Context, orderID string) (*orders.Order, error) {
	// First get the order
	order, err := r.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Then get the order items
	items, err := r.FindItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	order.Items = items
	return order, nil
}

// FindItemsByOrderID finds all items for a specific order
func (r *OrderRepository) FindItemsByOrderID(ctx context.Context, orderID string) ([]orders.OrderItem, error) {
	query := `
		SELECT item_id, order_id, product_name, quantity, created_at, updated_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []orders.OrderItem
	for rows.Next() {
		var item orders.OrderItem
		if err := rows.Scan(&item.ItemID, &item.OrderID, &item.ProductName, &item.Quantity, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return items, nil
}

// FindByUserID finds all orders for a specific user
func (r *OrderRepository) FindByUserID(ctx context.Context, userID string) ([]*orders.Order, error) {
	query := `
		SELECT order_id, order_number, created_by, assigned_to, delivery_address, status, created_at, updated_at
		FROM orders
		WHERE created_by = $1
		ORDER BY created_at DESC
	`
	return r.findOrdersWithQuery(ctx, query, userID)
}

// FindByAssignedTo finds all orders assigned to a specific user
func (r *OrderRepository) FindByAssignedTo(ctx context.Context, userID string) ([]*orders.Order, error) {
	query := `
		SELECT order_id, order_number, created_by, assigned_to, delivery_address, status, created_at, updated_at
		FROM orders
		WHERE assigned_to = $1
		ORDER BY created_at DESC
	`
	return r.findOrdersWithQuery(ctx, query, userID)
}

// FindByStatus finds all orders with a specific status
func (r *OrderRepository) FindByStatus(ctx context.Context, status string) ([]*orders.Order, error) {
	query := `
		SELECT order_id, order_number, created_by, assigned_to, delivery_address, status, created_at, updated_at
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
		SELECT order_id, order_number, created_by, assigned_to, delivery_address, status, created_at, updated_at
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
		if err := rows.Scan(&o.OrderID, &o.OrderNumber, &o.CreatedBy, &o.AssignedTo, &o.DeliveryAddress, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		ordersList = append(ordersList, &o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return ordersList, nil
}

// ============ Order Items Repository Methods ============

// SaveOrderItem saves a single order item to the database
func (r *OrderRepository) SaveOrderItem(ctx context.Context, item *orders.OrderItem) error {
	query := `
		INSERT INTO order_items (item_id, order_id, product_name, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Pool.Exec(ctx, query, item.ItemID, item.OrderID, item.ProductName, item.Quantity, item.CreatedAt, item.UpdatedAt)
	return err
}

// SaveOrderItems saves multiple order items to the database in a single transaction
func (r *OrderRepository) SaveOrderItems(ctx context.Context, items []orders.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO order_items (item_id, order_id, product_name, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, item := range items {
		_, err := tx.Exec(ctx, query, item.ItemID, item.OrderID, item.ProductName, item.Quantity, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// FindOrderItemByID finds a specific order item by its ID
func (r *OrderRepository) FindOrderItemByID(ctx context.Context, itemID string) (*orders.OrderItem, error) {
	query := `
		SELECT item_id, order_id, product_name, quantity, created_at, updated_at
		FROM order_items
		WHERE item_id = $1
	`
	var item orders.OrderItem
	err := r.db.Pool.QueryRow(ctx, query, itemID).Scan(
		&item.ItemID, &item.OrderID, &item.ProductName, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateOrderItem updates an existing order item
func (r *OrderRepository) UpdateOrderItem(ctx context.Context, item *orders.OrderItem) error {
	query := `
		UPDATE order_items 
		SET product_name = $2, quantity = $3, updated_at = $4
		WHERE item_id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, item.ItemID, item.ProductName, item.Quantity, item.UpdatedAt)
	return err
}

// DeleteOrderItem deletes a specific order item
func (r *OrderRepository) DeleteOrderItem(ctx context.Context, itemID string) error {
	query := `DELETE FROM order_items WHERE item_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, itemID)
	return err
}

// DeleteOrderItemsByOrderID deletes all items for a specific order
func (r *OrderRepository) DeleteOrderItemsByOrderID(ctx context.Context, orderID string) error {
	query := `DELETE FROM order_items WHERE order_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, orderID)
	return err
}
