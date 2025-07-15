package user

import "context"

// Repository defines the interface for user data storage.
type Repository interface {
	// Basic CRUD operations
	Save(ctx context.Context, user *User, userData *UserData) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, userID string) (*User, *UserData, error)
	Update(ctx context.Context, user *User, userData *UserData) error
	UpdateLastConnection(ctx context.Context, userID string) error
	Delete(ctx context.Context, userID string) error

	// Query operations
	FindByRole(ctx context.Context, role string) ([]*User, error)
	FindByRoleWithData(ctx context.Context, role string) ([]*UserWithData, error)
	FindAll(ctx context.Context, limit, offset int) ([]*User, int, error)
	FindAllWithData(ctx context.Context, limit, offset int) ([]*UserWithData, int, error)

	// Utility operations
	Exists(ctx context.Context, userID string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// Location operations
	UpdateLocation(ctx context.Context, userID string, latitude, longitude float64) error
	GetLocation(ctx context.Context, userID string) (*UserLocation, error)
	GetActiveDriversWithLocations(ctx context.Context) ([]*DriverLocation, error)
}
