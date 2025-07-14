package orders

import "context"

// Repository defines the interface for order data storage.
type Repository interface {
	// Basic CRUD operations
	Save(ctx context.Context, o *Order) error
	FindByID(ctx context.Context, orderID string) (*Order, error)
	FindByIDWithItems(ctx context.Context, orderID string) (*Order, error)
	FindAll(ctx context.Context) ([]*Order, error)
	Update(ctx context.Context, o *Order) error
	Delete(ctx context.Context, orderID string) error

	// Query operations
	FindByUserID(ctx context.Context, userID string) ([]*Order, error)
	FindByStatus(ctx context.Context, status string) ([]*Order, error)

	// Pagination
	FindWithPagination(ctx context.Context, limit, offset int) ([]*Order, int, error)

	// Utility operations
	Exists(ctx context.Context, orderID string) (bool, error)
	CountByStatus(ctx context.Context, status string) (int, error)

	// Order items operations
	FindItemsByOrderID(ctx context.Context, orderID string) ([]OrderItem, error)
	SaveOrderItem(ctx context.Context, item *OrderItem) error
	SaveOrderItems(ctx context.Context, items []OrderItem) error
	FindOrderItemByID(ctx context.Context, itemID string) (*OrderItem, error)
	UpdateOrderItem(ctx context.Context, item *OrderItem) error
	DeleteOrderItem(ctx context.Context, itemID string) error
	DeleteOrderItemsByOrderID(ctx context.Context, orderID string) error
}
