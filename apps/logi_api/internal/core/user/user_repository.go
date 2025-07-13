package user

import "context"

// Repository defines the interface for user data storage.
type Repository interface {
	Save(ctx context.Context, user *User, userData *UserData) error
}
