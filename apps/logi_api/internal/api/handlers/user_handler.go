package handlers

import (
	"bombayv/logiapp-monorepo/logi_api/internal/core/user"
	"fmt"
	"net/http"

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
	Phone     string `json:"phone" binding:"required"`
	Role      string `json:"role" binding:"required"`
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
	// Remove "Bearer " prefix
	tokenString = tokenString[len("Bearer "):]

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
		"user": user.Email,
		"role": user.Role,
		"profile": gin.H{
			"first_name":   userData.FirstName,
			"last_name":    userData.LastName,
			"phone_number": userData.PhoneNumber,
		},
	})
}
