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
func (s *Service) CreateOrder(ctx context.Context, email, address string) (*Order, error) {
	orderID := uuid.New().String()
	now := time.Now()

	// Find the user by email to get the user ID
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	order := &Order{
		OrderID:         orderID,
		CreatedBy:       user.UserID,
		AssignedTo:      "",
		DeliveryAddress: address,
		Status:          "pending",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}

// FindAll retrieves all orders from the repository.
func (s *Service) FindAll(ctx context.Context) ([]*Order, error) {
	orders, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve orders: %w", err)
	}
	return orders, nil
}

// FindByID retrieves an order by its ID.
func (s *Service) FindByID(ctx context.Context, orderID string) (*Order, error) {
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
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
	return orders, nil
}

// FindByStatus retrieves all orders with a specific status.
func (s *Service) FindByStatus(ctx context.Context, status string) ([]*Order, error) {
	orders, err := s.repo.FindByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve orders by status: %w", err)
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
		order.AssignedTo = assignedTo
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
