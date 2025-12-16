package orders

import (
	"bombayv/logiapp-monorepo/logi_api/internal/config"
	"bombayv/logiapp-monorepo/logi_api/internal/core/user"
	"bombayv/logiapp-monorepo/logi_api/internal/email"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/cache"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service provides order-related operations.
type Service struct {
	repo         Repository
	userRepo     user.Repository
	cache        *cache.Cache
	emailService *email.Service
}

// NewService creates a new order service.
func NewService(repo Repository, userRepo user.Repository, cache *cache.Cache, emailService *email.Service) *Service {
	return &Service{repo: repo, userRepo: userRepo, cache: cache, emailService: emailService}
}

// CreateOrder handles the business logic for creating a new order.
func (s *Service) CreateOrder(ctx context.Context, email, orderName, orderPhoneNumber string, orderEmail, orderCedula, address string, orderNumber string) (*Order, error) {
	orderID := uuid.New().String()
	now := time.Now()

	// Find the user by email to get the user ID
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Convert empty strings to nil for optional fields
	var orderEmailPtr *string
	if orderEmail != "" {
		orderEmailPtr = &orderEmail
	}

	var orderCedulaPtr *string
	if orderCedula != "" {
		orderCedulaPtr = &orderCedula
	}

	order := &Order{
		OrderID:          orderID,
		OrderNumber:      orderNumber,
		CreatedBy:        user.UserID,
		AssignedTo:       nil,
		OrderName:        orderName,
		OrderPhoneNumber: orderPhoneNumber,
		OrderEmail:       orderEmailPtr,
		OrderCedula:      orderCedulaPtr,
		DeliveryAddress:  address,
		Status:           "pending",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// Populate username fields
	if err := s.populateUsernames(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to populate usernames: %w", err)
	}

	// Invalidate orders cache since we added a new order
	s.invalidateOrderCache(ctx, order.OrderID)

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
	// Try to get from cache first
	cacheKey := "orders_all"

	var cached []*Order
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	// If not in cache, get from database
	orders, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve orders: %w", err)
	}

	// Populate username fields for all orders
	if err := s.populateUsernamesForOrders(ctx, orders); err != nil {
		return nil, fmt.Errorf("failed to populate usernames: %w", err)
	}

	// Cache the result for 1 minute (short TTL since orders change frequently)
	if cacheErr := s.cache.Set(ctx, cacheKey, orders, 1*time.Minute); cacheErr != nil {
		// Log cache error but don't fail the request
		fmt.Printf("Failed to cache all orders: %v\n", cacheErr)
	}

	return orders, nil
}

// FindByID retrieves an order by its ID.
func (s *Service) FindByID(ctx context.Context, orderID string) (*Order, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("order:%s", orderID)

	var cached Order
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	// If not in cache, get from database
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	// Populate username fields
	if err := s.populateUsernames(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to populate usernames: %w", err)
	}

	// Cache the result for 3 minutes
	if cacheErr := s.cache.Set(ctx, cacheKey, order, 3*time.Minute); cacheErr != nil {
		// Log cache error but don't fail the request
		fmt.Printf("Failed to cache order %s: %v\n", orderID, cacheErr)
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

// FindByAssignedTo retrieves pending and in_progress orders assigned to a specific user (driver).
func (s *Service) FindByAssignedTo(ctx context.Context, userID string) ([]*Order, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("orders_assigned_to:%s", userID)

	var cached []*Order
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	// Validate that the user exists
	exists, err := s.userRepo.Exists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// If not in cache, get from database
	orders, err := s.repo.FindByAssignedTo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve assigned orders: %w", err)
	}

	// Populate username fields for all orders
	if err := s.populateUsernamesForOrders(ctx, orders); err != nil {
		return nil, fmt.Errorf("failed to populate usernames: %w", err)
	}

	// Cache the result for 1 minute (short TTL since orders change frequently)
	if cacheErr := s.cache.Set(ctx, cacheKey, orders, 1*time.Minute); cacheErr != nil {
		// Log cache error but don't fail the request
		fmt.Printf("Failed to cache assigned orders for user %s: %v\n", userID, cacheErr)
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

	// Try to get from cache first
	cacheKey := fmt.Sprintf("orders_paginated:%d:%d", limit, offset)

	type OrdersWithPaginationResult struct {
		Orders []*Order `json:"orders"`
		Total  int      `json:"total"`
	}

	var cached OrdersWithPaginationResult
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached.Orders, cached.Total, nil
	}

	// If not in cache, get from database
	orders, total, err := s.repo.FindWithPagination(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve paginated orders: %w", err)
	}

	// Populate username fields for all orders
	if err := s.populateUsernamesForOrders(ctx, orders); err != nil {
		return nil, 0, fmt.Errorf("failed to populate usernames: %w", err)
	}

	// Cache the result for 1 minute (short TTL since orders change frequently)
	result := OrdersWithPaginationResult{Orders: orders, Total: total}
	if cacheErr := s.cache.Set(ctx, cacheKey, result, 1*time.Minute); cacheErr != nil {
		// Log cache error but don't fail the request
		fmt.Printf("Failed to cache paginated orders (limit:%d, offset:%d): %v\n", limit, offset, cacheErr)
	}

	return orders, total, nil
}

// UpdateOrder updates an existing order.
func (s *Service) UpdateOrder(ctx context.Context, orderID string, assignedTo string, assignedToProvided bool, address, status string) (*Order, error) {
	// Find the existing order
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Update fields
	if assignedToProvided {
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
		} else {
			// Set to null when assignedTo is empty string but was provided
			order.AssignedTo = nil
		}
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

	// Invalidate orders cache since we updated an order
	s.invalidateOrderCache(ctx, order.OrderID)

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

	// Invalidate orders cache since we deleted an order
	s.invalidateOrderCache(ctx, orderID)

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

// invalidateOrderCache invalidates all cache entries related to orders
func (s *Service) invalidateOrderCache(ctx context.Context, orderID string) {
	// Invalidate specific order cache
	orderKey := fmt.Sprintf("order:%s", orderID)
	if err := s.cache.Delete(ctx, orderKey); err != nil {
		fmt.Printf("Failed to invalidate order cache for %s: %v\n", orderID, err)
	}

	// Invalidate all orders cache
	if err := s.cache.Delete(ctx, "orders_all"); err != nil {
		fmt.Printf("Failed to invalidate all orders cache: %v\n", err)
	}

	// Invalidate paginated orders cache (all combinations of limit/offset)
	if err := s.cache.DeletePattern(ctx, "orders_paginated:*"); err != nil {
		fmt.Printf("Failed to invalidate paginated orders cache: %v\n", err)
	}

	// Invalidate assigned orders cache (all users)
	if err := s.cache.DeletePattern(ctx, "orders_assigned_to:*"); err != nil {
		fmt.Printf("Failed to invalidate assigned orders cache: %v\n", err)
	}
}

// generatePublicID generates a random 6-character alphanumeric string
func generatePublicID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // Ensure different seed for each char if called rapidly
	}
	return string(b)
}

// CreateOrderForm creates a new satisfaction form for an order
func (s *Service) CreateOrderForm(ctx context.Context, orderID string, driverID *string, driverRating *int, cargoCondition, comments *string) (*OrderForm, error) {
	// Check if a form already exists for this order
	existingForm, err := s.repo.FindOrderFormByOrderID(ctx, orderID)
	if err == nil && existingForm != nil {
		// If we are just requesting to send the survey (no rating provided) and an unfinished form exists, reuse it
		if driverRating == nil && !existingForm.IsFinished {
			// Resend email asynchronously
			go func() {
				bgCtx := context.Background()
				order, err := s.repo.FindByID(bgCtx, orderID)
				if err != nil {
					fmt.Printf("failed to fetch order for email sending: %v\n", err)
					return
				}

				if order.OrderEmail != nil && *order.OrderEmail != "" {
					baseURL := config.App.WebBaseURL
					if baseURL == "" {
						baseURL = "http://localhost:5173"
					}
					surveyLink := fmt.Sprintf("%s/encuestas/%s", baseURL, existingForm.PublicID)

					err := s.emailService.SendSurveyEmail([]string{*order.OrderEmail}, surveyLink)
					if err != nil {
						fmt.Printf("failed to send survey email: %v\n", err)
					} else {
						fmt.Printf("survey email sent to %s\n", *order.OrderEmail)
					}
				}
			}()
			return existingForm, nil
		}
	}

	formID := uuid.New().String()
	publicID := generatePublicID()
	now := time.Now()

	isFinished := false
	if driverRating != nil {
		isFinished = true
	}

	form := &OrderForm{
		FormID:         formID,
		PublicID:       publicID,
		OrderID:        orderID,
		DriverID:       driverID,
		DriverRating:   driverRating,
		CargoCondition: cargoCondition,
		Comments:       comments,
		IsFinished:     isFinished,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.SaveOrderForm(ctx, form); err != nil {
		return nil, fmt.Errorf("failed to save order form: %w", err)
	}

	// Send email asynchronously
	go func() {
		// Create a new context for the background task
		bgCtx := context.Background()

		order, err := s.repo.FindByID(bgCtx, orderID)
		if err != nil {
			fmt.Printf("failed to fetch order for email sending: %v\n", err)
			return
		}

		if order.OrderEmail != nil && *order.OrderEmail != "" {
			baseURL := config.App.WebBaseURL
			if baseURL == "" {
				baseURL = "http://localhost:5173"
			}
			surveyLink := fmt.Sprintf("%s/encuestas/%s", baseURL, publicID)

			err := s.emailService.SendSurveyEmail([]string{*order.OrderEmail}, surveyLink)
			if err != nil {
				fmt.Printf("failed to send survey email: %v\n", err)
			} else {
				fmt.Printf("survey email sent to %s\n", *order.OrderEmail)
			}
		}
	}()

	return form, nil
}

// GetOrderForm retrieves a satisfaction form for a specific order
func (s *Service) GetOrderForm(ctx context.Context, orderID string) (*OrderForm, error) {
	form, err := s.repo.FindOrderFormByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order form: %w", err)
	}

	return form, nil
}

// GetOrderFormByPublicID retrieves a satisfaction form by its public ID
func (s *Service) GetOrderFormByPublicID(ctx context.Context, publicID string) (*OrderForm, error) {
	form, err := s.repo.FindOrderFormByPublicID(ctx, publicID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order form: %w", err)
	}

	// Populate driver info if available
	if form.DriverID != nil {
		driver, driverData, err := s.userRepo.FindByID(ctx, *form.DriverID)
		if err == nil {
			if driverData != nil {
				form.DriverName = driverData.FirstName + " " + driverData.LastName
			}
			if driver != nil {
				form.DriverEmail = driver.Email
			}
		}
	}

	return form, nil
}

// UpdateOrderForm updates an existing satisfaction form
func (s *Service) UpdateOrderForm(ctx context.Context, formID string, driverRating *int, cargoCondition, comments *string) (*OrderForm, error) {
	form, err := s.repo.FindOrderFormByID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("order form not found: %w", err)
	}

	// Update fields
	if driverRating != nil {
		form.DriverRating = driverRating
		form.IsFinished = true
	}

	if cargoCondition != nil {
		form.CargoCondition = cargoCondition
	}

	if comments != nil {
		form.Comments = comments
	}

	form.UpdatedAt = time.Now()

	if err := s.repo.UpdateOrderForm(ctx, form); err != nil {
		return nil, fmt.Errorf("failed to update order form: %w", err)
	}

	return form, nil
}

// SubmitOrderForm submits a satisfaction form using its public ID
func (s *Service) SubmitOrderForm(ctx context.Context, publicID string, driverRating *int, cargoCondition, comments *string) (*OrderForm, error) {
	form, err := s.repo.FindOrderFormByPublicID(ctx, publicID)
	if err != nil {
		return nil, fmt.Errorf("order form not found: %w", err)
	}

	if form.IsFinished {
		return nil, fmt.Errorf("form already submitted")
	}

	// Update fields
	if driverRating != nil {
		form.DriverRating = driverRating
	}

	if cargoCondition != nil {
		form.CargoCondition = cargoCondition
	}

	if comments != nil {
		form.Comments = comments
	}

	form.IsFinished = true
	form.UpdatedAt = time.Now()

	if err := s.repo.UpdateOrderForm(ctx, form); err != nil {
		return nil, fmt.Errorf("failed to update order form: %w", err)
	}

	return form, nil
}

// DeleteOrderForm deletes a satisfaction form
func (s *Service) DeleteOrderForm(ctx context.Context, formID string) error {
	// Check if form exists
	_, err := s.repo.FindOrderFormByID(ctx, formID)
	if err != nil {
		return fmt.Errorf("order form not found: %w", err)
	}

	if err := s.repo.DeleteOrderForm(ctx, formID); err != nil {
		return fmt.Errorf("failed to delete order form: %w", err)
	}

	return nil
}
