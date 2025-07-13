package handlers

import (
	"bombayv/logiapp-monorepo/logi_api/internal/core/user"
	"bombayv/logiapp-monorepo/logi_api/internal/utils"
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
func LoginUser(c *gin.Context) {
	// In a real app, you would:
	// 1. Bind the request body to a login credentials struct.
	// 2. Validate the credentials with the user service.
	// 3. If valid, generate a JWT.
	// 4. Return the token.
	c.JSON(http.StatusOK, gin.H{"token": "fake-jwt-token-for-testing"})
}

func GetToken(c *gin.Context) {
	// In a real app, you would generate a JWT token here.
	// For now, we return a fake token for testing purposes.
	token, err := utils.GenerateJWT("test-user-id", "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Token generated successfully",
	})
}
