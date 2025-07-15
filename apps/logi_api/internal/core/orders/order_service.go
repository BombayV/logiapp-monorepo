package orders

import (
	"bombayv/logiapp-monorepo/logi_api/internal/core/user"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service provides order-related operations.
type Service struct {
	repo     Repository
	userRepo user.Repository
}

// NewService creates a new order service.
func NewService(repo Repository, userRepo user.Repository) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

// CreateOrder handles the business logic for creating a new order.
func (s *Service) CreateOrder(ctx context.Context, email, address string, orderNumber string) (*Order, error) {
	orderID := uuid.New().String()
	now := time.Now()

	// Find the user by email to get the user ID
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	order := &Order{
		OrderID:         orderID,
		OrderNumber:     orderNumber,
		CreatedBy:       user.UserID,
		AssignedTo:      nil,
		DeliveryAddress: address,
		Status:          "pending",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// Populate username fields
	if err := s.populateUsernames(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to populate usernames: %w", err)
	}

	return order, nil
}

// populateUsernames populates the username fields for a single order
func (s *Service) populateUsernames(ctx context.Context, order *Order) error {
	// Get created_by username
	createdByUser, _, err := s.userRepo.FindByID(ctx, order.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to find created_by user: %w", err)
	}
	order.CreatedByUsername = createdByUser.Email

	// Get assigned_to username if exists
	if order.AssignedTo != nil {
		assignedToUser, _, err := s.userRepo.FindByID(ctx, *order.AssignedTo)
		if err != nil {
			return fmt.Errorf("failed to find assigned_to user: %w", err)
		}
		order.AssignedToUsername = &assignedToUser.Email
	}

	return nil
}

// populateUsernamesForOrders populates the username fields for multiple orders
func (s *Service) populateUsernamesForOrders(ctx context.Context, orders []*Order) error {
	for _, order := range orders {
		if err := s.populateUsernames(ctx, order); err != nil {
			return err
		}
	}
	return nil
}

// FindAll retrieves all orders from the repository.
func (s *Service) FindAll(ctx context.Context) ([]*Order, error) {
	orders, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve orders: %w", err)
	}

	// Populate username fields for all orders
	if err := s.populateUsernamesForOrders(ctx, orders); err != nil {
		return nil, fmt.Errorf("failed to populate usernames: %w", err)
	}

	return orders, nil
}

// FindByID retrieves an order by its ID.
func (s *Service) FindByID(ctx context.Context, orderID string) (*Order, error) {
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	// Populate username fields
	if err := s.populateUsernames(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to populate usernames: %w", err)
	}

	return order, nil
}

// FindByIDWithItems retrieves an order by its ID including all order items.
func (s *Service) FindByIDWithItems(ctx context.Context, orderID string) (*Order, error) {
	order, err := s.repo.FindByIDWithItems(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order with items: %w", err)
	}

	// Populate username fields
	if err := s.populateUsernames(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to populate usernames: %w", err)
	}

	return order, nil
}

// FindByUserID retrieves all orders for a specific user.
func (s *Service) FindByUserID(ctx context.Context, userID string) ([]*Order, error) {
	// Validate that the user exists
	exists, err := s.userRepo.Exists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	orders, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user orders: %w", err)
	}

	// Populate username fields for all orders
	if err := s.populateUsernamesForOrders(ctx, orders); err != nil {
		return nil, fmt.Errorf("failed to populate usernames: %w", err)
	}

	return orders, nil
}

// FindByStatus retrieves all orders with a specific status.
func (s *Service) FindByStatus(ctx context.Context, status string) ([]*Order, error) {
	orders, err := s.repo.FindByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve orders by status: %w", err)
	}

	// Populate username fields for all orders
	if err := s.populateUsernamesForOrders(ctx, orders); err != nil {
		return nil, fmt.Errorf("failed to populate usernames: %w", err)
	}

	return orders, nil
}

// FindWithPagination retrieves orders with pagination.
func (s *Service) FindWithPagination(ctx context.Context, limit, offset int) ([]*Order, int, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}
	if offset < 0 {
		offset = 0
	}

	orders, total, err := s.repo.FindWithPagination(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve paginated orders: %w", err)
	}

	// Populate username fields for all orders
	if err := s.populateUsernamesForOrders(ctx, orders); err != nil {
		return nil, 0, fmt.Errorf("failed to populate usernames: %w", err)
	}

	return orders, total, nil
}

// UpdateOrder updates an existing order.
func (s *Service) UpdateOrder(ctx context.Context, orderID string, assignedTo, address, status string) (*Order, error) {
	// Find the existing order
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Update fields
	if assignedTo != "" {
		// Validate that the assigned user exists
		exists, err := s.userRepo.Exists(ctx, assignedTo)
		if err != nil {
			return nil, fmt.Errorf("failed to validate assigned user: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("assigned user not found")
		}
		order.AssignedTo = &assignedTo
	}

	if address != "" {
		order.DeliveryAddress = address
	}

	if status != "" {
		order.Status = status
	}

	order.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Populate username fields
	if err := s.populateUsernames(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to populate usernames: %w", err)
	}

	return order, nil
}

// DeleteOrder deletes an order.
func (s *Service) DeleteOrder(ctx context.Context, orderID string) error {
	// Check if order exists
	exists, err := s.repo.Exists(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to check order existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("order not found")
	}

	if err := s.repo.Delete(ctx, orderID); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}

// GetOrderStats returns statistics about orders.
func (s *Service) GetOrderStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)

	statuses := []string{"pending", "assigned", "in_progress", "completed", "cancelled"}

	for _, status := range statuses {
		count, err := s.repo.CountByStatus(ctx, status)
		if err != nil {
			return nil, fmt.Errorf("failed to get stats for status %s: %w", status, err)
		}
		stats[status] = count
	}

	return stats, nil
}

// ============ Order Items Service Methods ============

// AddOrderItem adds a single item to an order.
func (s *Service) AddOrderItem(ctx context.Context, orderID, productName string, quantity int) (*OrderItem, error) {
	// Validate that the order exists
	exists, err := s.repo.Exists(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check order existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("order not found")
	}

	// Validate quantity
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	itemID := uuid.New().String()
	now := time.Now()

	item := &OrderItem{
		ItemID:      itemID,
		OrderID:     orderID,
		ProductName: productName,
		Quantity:    quantity,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.SaveOrderItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save order item: %w", err)
	}

	return item, nil
}

// AddOrderItems adds multiple items to an order in a single operation.
func (s *Service) AddOrderItems(ctx context.Context, orderID string, itemRequests []struct {
	ProductName string `json:"product_name" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
}) ([]OrderItem, error) {
	// Validate that the order exists
	exists, err := s.repo.Exists(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check order existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("order not found")
	}

	if len(itemRequests) == 0 {
		return nil, fmt.Errorf("no items provided")
	}

	now := time.Now()
	items := make([]OrderItem, len(itemRequests))

	for i, req := range itemRequests {
		if req.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0 for item %d", i+1)
		}

		items[i] = OrderItem{
			ItemID:      uuid.New().String(),
			OrderID:     orderID,
			ProductName: req.ProductName,
			Quantity:    req.Quantity,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}

	if err := s.repo.SaveOrderItems(ctx, items); err != nil {
		return nil, fmt.Errorf("failed to save order items: %w", err)
	}

	return items, nil
}

// GetOrderItems retrieves all items for a specific order.
func (s *Service) GetOrderItems(ctx context.Context, orderID string) ([]OrderItem, error) {
	// Validate that the order exists
	exists, err := s.repo.Exists(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check order existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("order not found")
	}

	items, err := s.repo.FindItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order items: %w", err)
	}

	return items, nil
}

// GetOrderItemByID retrieves a specific order item by its ID.
func (s *Service) GetOrderItemByID(ctx context.Context, itemID string) (*OrderItem, error) {
	item, err := s.repo.FindOrderItemByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order item: %w", err)
	}
	return item, nil
}

// UpdateOrderItem updates an existing order item.
func (s *Service) UpdateOrderItem(ctx context.Context, itemID, productName string, quantity int) (*OrderItem, error) {
	// Find the existing item
	item, err := s.repo.FindOrderItemByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("order item not found: %w", err)
	}

	// Update fields
	if productName != "" {
		item.ProductName = productName
	}

	if quantity > 0 {
		item.Quantity = quantity
	}

	item.UpdatedAt = time.Now()

	if err := s.repo.UpdateOrderItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to update order item: %w", err)
	}

	return item, nil
}

// DeleteOrderItem deletes a specific order item.
func (s *Service) DeleteOrderItem(ctx context.Context, itemID string) error {
	// Check if item exists
	_, err := s.repo.FindOrderItemByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("order item not found: %w", err)
	}

	if err := s.repo.DeleteOrderItem(ctx, itemID); err != nil {
		return fmt.Errorf("failed to delete order item: %w", err)
	}

	return nil
}
