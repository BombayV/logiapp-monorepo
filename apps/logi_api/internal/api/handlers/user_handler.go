package handlers

import (
	"bombayv/logiapp-monorepo/logi_api/internal/core/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler holds the dependencies for the user handlers.
type UserHandler struct {
	UserService *user.Service
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{UserService: userService}
}

type RegisterUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone_number" binding:"required"`
	Role      string `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Phone     *string `json:"phone_number,omitempty"`
	Role      *string `json:"role,omitempty"`
}

// RegisterUser handles new user registration.
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var request RegisterUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		fmt.Println("Error binding request body:", err)
		return
	}

	// Validate role against allowed values from schema
	validRoles := map[string]bool{
		"admin":  true,
		"driver": true,
		"sales":  true,
	}

	if !validRoles[request.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Allowed values: admin, driver, sales"})
		return
	}

	newUser, err := h.UserService.CreateUser(
		c.Request.Context(),
		request.Email,
		request.Password,
		request.FirstName,
		request.LastName,
		request.Phone,
		request.Role,
	)

	if err != nil {
		switch err {
		case user.ErrEmailExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case user.ErrInvalidCredentials, user.ErrPasswordTooWeak:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user":    newUser,
	})
}

// LoginUser handles user authentication.
func (h *UserHandler) Login(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := h.UserService.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		switch err {
		case user.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout handles user logout.
func (h *UserHandler) Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	// Remove "Bearer " prefix
	const bearerPrefix = "Bearer "
	if len(tokenString) < len(bearerPrefix) || tokenString[:len(bearerPrefix)] != bearerPrefix {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization header format"})
		return
	}
	tokenString = tokenString[len(bearerPrefix):]

	err := h.UserService.Logout(c.Request.Context(), tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// Me handles getting the current user's profile information.
func (h *UserHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	user, userData, err := h.UserService.GetUserProfile(c.Request.Context(), userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	// Return user profile without sensitive information
	c.JSON(http.StatusOK, gin.H{
		"user_id": user.UserID,
		"email":   user.Email,
		"role":    user.Role,
		"profile": gin.H{
			"first_name":      userData.FirstName,
			"last_name":       userData.LastName,
			"phone_number":    userData.PhoneNumber,
			"last_connection": userData.LastConnection,
		},
	})
}

// GetUsers returns all users with pagination (admin only)
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Check if user is admin (this would typically be done by middleware)
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

	users, total, err := h.UserService.GetAllUsersWithData(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format the response to include user data
	var formattedUsers []gin.H
	for _, userWithData := range users {
		userData := gin.H{
			"user_id":    userWithData.User.UserID,
			"email":      userWithData.User.Email,
			"role":       userWithData.User.Role,
			"created_at": userWithData.User.CreatedAt,
			"updated_at": userWithData.User.UpdatedAt,
		}

		// Add user data if available
		if userWithData.UserData != nil {
			userData["first_name"] = userWithData.UserData.FirstName
			userData["last_name"] = userWithData.UserData.LastName
			userData["phone_number"] = userWithData.UserData.PhoneNumber
			userData["last_connection"] = userWithData.UserData.LastConnection
		}

		formattedUsers = append(formattedUsers, userData)
	}

	c.JSON(http.StatusOK, gin.H{
		"users":  formattedUsers,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetUserByID returns a specific user by ID (admin only)
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Check if user is admin (this would typically be done by middleware)
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, userData, err := h.UserService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return full user profile for admin
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"user_id":    user.UserID,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
		"profile": gin.H{
			"first_name":      userData.FirstName,
			"last_name":       userData.LastName,
			"phone_number":    userData.PhoneNumber,
			"last_connection": userData.LastConnection,
		},
	})
}

// DeleteUser deletes a user (admin only)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	err := h.UserService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UpdateLocation handles location updates for drivers
func (h *UserHandler) UpdateLocation(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var request struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.UserService.UpdateLocation(c.Request.Context(), userIDStr, request.Latitude, request.Longitude)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location updated successfully"})
}

// GetLocation handles location retrieval for drivers
func (h *UserHandler) GetLocation(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	location, err := h.UserService.GetLocation(c.Request.Context(), userIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"latitude":   location.Latitude,
		"longitude":  location.Longitude,
		"updated_at": location.UpdatedAt,
	})
}

// GetActiveDriversWithLocations returns all drivers who have been active in the last 10 minutes with their locations
func (h *UserHandler) GetActiveDriversWithLocations(c *gin.Context) {
	drivers, err := h.UserService.GetActiveDriversWithLocations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve active drivers"})
		return
	}

	// Ensure drivers is always an empty array instead of null
	if drivers == nil {
		drivers = []*user.DriverLocation{}
	}

	c.JSON(http.StatusOK, gin.H{
		"drivers": drivers,
		"count":   len(drivers),
	})
}

// GetAllDrivers returns all drivers ordered by last connection
func (h *UserHandler) GetAllDrivers(c *gin.Context) {
	drivers, err := h.UserService.GetAllDrivers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve drivers"})
		return
	}

	// Ensure drivers is always an empty array instead of null
	if drivers == nil {
		drivers = []*user.Driver{}
	}

	c.JSON(http.StatusOK, gin.H{
		"drivers": drivers,
		"count":   len(drivers),
	})
}

type ResetPasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ResetPassword handles password reset requests
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var request ResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Reset password
	err := h.UserService.ResetPassword(c.Request.Context(), userIDStr, request.CurrentPassword, request.NewPassword)
	if err != nil {
		switch err {
		case user.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		case user.ErrPasswordTooWeak:
			c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be at least 8 characters long"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
