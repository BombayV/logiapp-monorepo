package user

import "context"

// Repository defines the interface for user data storage.
type Repository interface {
	// Basic CRUD operations
	Save(ctx context.Context, user *User, userData *UserData) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, userID string) (*User, *UserData, error)
	Update(ctx context.Context, user *User, userData *UserData) error
	Delete(ctx context.Context, userID string) error

	// Query operations
	FindByRole(ctx context.Context, role string) ([]*User, error)
	FindAll(ctx context.Context, limit, offset int) ([]*User, int, error)

	// Utility operations
	Exists(ctx context.Context, userID string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
