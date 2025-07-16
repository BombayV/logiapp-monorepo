package handlers

import (
	ordersCore "bombayv/logiapp-monorepo/logi_api/internal/core/orders"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// OrderHandler handles HTTP requests for orders.
type OrderHandler struct {
	service *ordersCore.Service
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(service *ordersCore.Service) *OrderHandler {
	return &OrderHandler{service: service}
}

type CreateOrderRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Address     string `json:"address" binding:"required"`
	OrderNumber string `json:"order_number" binding:"required,min=1,max=6"`
}

type UpdateOrderRequest struct {
	AssignedTo *string `json:"assigned_to,omitempty"`
	Address    *string `json:"address,omitempty"`
	Status     *string `json:"status,omitempty"`
}

type AddOrderItemRequest struct {
	ProductName string `json:"product_name" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
}

type AddOrderItemsRequest struct {
	Items []AddOrderItemRequest `json:"items" binding:"required,min=1"`
}

type UpdateOrderItemRequest struct {
	ProductName *string `json:"product_name,omitempty"`
	Quantity    *int    `json:"quantity,omitempty" binding:"omitempty,min=1"`
}

// CreateOrder handles the creation of a new order.
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate order_number: must be numeric only
	if !regexp.MustCompile(`^[0-9]+$`).MatchString(req.OrderNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_number must contain only numbers"})
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), req.Email, req.Address, req.OrderNumber)
	if err != nil {
		// Check if it's a duplicate order number error
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") &&
			strings.Contains(err.Error(), "order_number") {
			c.JSON(http.StatusConflict, gin.H{"error": "Ya existe una orden con este número"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetAllOrders returns all orders with optional pagination
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	// Check if pagination parameters are provided
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	if limitStr != "" || offsetStr != "" {
		h.GetOrdersWithPagination(c)
		return
	}

	ordersList, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure orders is always an empty array instead of null
	if ordersList == nil {
		ordersList = []*ordersCore.Order{}
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": ordersList,
		"total":  len(ordersList),
	})
}

// GetOrdersWithPagination returns paginated orders
func (h *OrderHandler) GetOrdersWithPagination(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	orders, total, err := h.service.FindWithPagination(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure orders is always an empty array instead of null
	if orders == nil {
		orders = []*ordersCore.Order{}
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetOrderByID returns a specific order by ID
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	order, err := h.service.FindByIDWithItems(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrder updates an existing order
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	// Read the request body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Parse as structured request
	var req UpdateOrderRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse as raw JSON to check field presence
	var rawJSON map[string]interface{}
	if err := json.Unmarshal(body, &rawJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debug logging
	fmt.Printf("Received update request: %+v\n", req)
	fmt.Printf("Raw JSON: %+v\n", rawJSON)

	var assignedToProvided bool
	var assignedToValue string

	// Check if assigned_to field exists in the JSON (even if null)
	if assignedToRaw, exists := rawJSON["assigned_to"]; exists {
		assignedToProvided = true
		if assignedToRaw != nil {
			assignedToValue = assignedToRaw.(string)
		}
		// If assignedToRaw is nil, assignedToValue stays empty string (for null assignment)
	}

	fmt.Printf("AssignedTo provided: %t, value: '%s'\n", assignedToProvided, assignedToValue)

	// Validate status if provided
	if req.Status != nil {
		validStatuses := map[string]bool{
			"pending":     true,
			"in_progress": true,
			"completed":   true,
			"cancelled":   true,
		}

		if !validStatuses[*req.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Allowed values: pending, in_progress, completed, cancelled"})
			return
		}
	}

	assignedTo := assignedToValue
	address := ""
	if req.Address != nil {
		address = *req.Address
	}

	status := ""
	if req.Status != nil {
		status = *req.Status
	}

	order, err := h.service.UpdateOrder(c.Request.Context(), orderID, assignedTo, assignedToProvided, address, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// DeleteOrder deletes an order
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	err := h.service.DeleteOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// GetOrderStats returns statistics about orders
func (h *OrderHandler) GetOrderStats(c *gin.Context) {
	stats, err := h.service.GetOrderStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// ============ Order Items Handlers ============

// AddOrderItem adds a single item to an order
func (h *OrderHandler) AddOrderItem(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	var req AddOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.AddOrderItem(c.Request.Context(), orderID, req.ProductName, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// AddOrderItems adds multiple items to an order
func (h *OrderHandler) AddOrderItems(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	var req AddOrderItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request items to service format
	itemRequests := make([]struct {
		ProductName string `json:"product_name" binding:"required"`
		Quantity    int    `json:"quantity" binding:"required,min=1"`
	}, len(req.Items))

	for i, item := range req.Items {
		itemRequests[i].ProductName = item.ProductName
		itemRequests[i].Quantity = item.Quantity
	}

	items, err := h.service.AddOrderItems(c.Request.Context(), orderID, itemRequests)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"items": items,
		"total": len(items),
	})
}

// GetOrderItems retrieves all items for a specific order
func (h *OrderHandler) GetOrderItems(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	items, err := h.service.GetOrderItems(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": len(items),
	})
}

// UpdateOrderItem updates a specific order item
func (h *OrderHandler) UpdateOrderItem(c *gin.Context) {
	itemID := c.Param("item_id")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item ID is required"})
		return
	}

	var req UpdateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productName := ""
	if req.ProductName != nil {
		productName = *req.ProductName
	}

	quantity := 0
	if req.Quantity != nil {
		quantity = *req.Quantity
	}

	item, err := h.service.UpdateOrderItem(c.Request.Context(), itemID, productName, quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteOrderItem deletes a specific order item
func (h *OrderHandler) DeleteOrderItem(c *gin.Context) {
	itemID := c.Param("item_id")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item ID is required"})
		return
	}

	err := h.service.DeleteOrderItem(c.Request.Context(), itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order item deleted successfully"})
}

// GetOrdersByUserID retrieves orders for a specific user (driver)
func (h *OrderHandler) GetOrdersByUserID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Get the authenticated user's ID from the JWT token
	authUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Drivers can only access their own orders
	if userID != authUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You can only view your own orders"})
		return
	}

	orders, err := h.service.FindByAssignedTo(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
