package user

import (
	"bombayv/logiapp-monorepo/logi_api/internal/utils"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrForbidden          = errors.New("forbidden: insufficient permissions")
)

// CreateUser handles the business logic for creating a new user.
func CreateUser(email, password, firstName, lastName, phone, role string) (*User, error) {
	fmt.Println("All parameters received in CreateUser:")
	fmt.Printf("Email: %s, Password: %s, FirstName: %s, LastName: %s, Phone: %s, Role: %s\n",
		email,
		password,
		firstName,
		lastName,
		phone,
		role,
	)
	// 2. Validate role
	if role != "sales" && role != "driver" {
		return nil, errors.New("invalid role provided: must be 'sales' or 'driver'")
	}

	// 3. Validate email and password strength.
	if email == "" || password == "" || len(password) < 8 {
		return nil, ErrInvalidCredentials
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
