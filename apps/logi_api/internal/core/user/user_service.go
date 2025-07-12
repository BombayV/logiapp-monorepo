package user

import (
	"bombayv/logiapp-monorepo/logi_api/internal/utils"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("missing or invalid credentials")
	ErrPasswordTooWeak    = errors.New("password must be at least 8 characters long")
	ErrForbidden          = errors.New("forbidden: insufficient permissions")
)

// CreateUser handles the business logic for creating a new user.
func CreateUser(email, password, firstName, lastName, phone, role string) (*User, error) {
	// 2. Validate role
	if role != "sales" && role != "driver" {
		return nil, errors.New("invalid role provided: must be 'sales' or 'driver'")
	}

	// 3. Validate email and password strength.
	if email == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	if len(password) < 8 {
		return nil, ErrPasswordTooWeak
	}

	// 4. Hash the password.
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 5. Create User and UserData structs.
	userID := uuid.New().String()
	now := time.Now()
	user := &User{
		UserID:       userID,
		Email:        email,
		PasswordHash: passwordHash,
		UserRole:     role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	fmt.Println("User created with password hash:", passwordHash)

	// userData := &UserData{
	// 	UserID:         userID,
	// 	FirstName:      firstName,
	// 	LastName:       lastName,
	// 	PhoneNumber:    phone,
	// 	LastConnection: now,
	// 	CreatedAt:      now,
	// 	UpdatedAt:      now,
	// }

	// 6. Call the repository to save the user to the database.

	// 7. Return the created user.
	return user, nil
}
