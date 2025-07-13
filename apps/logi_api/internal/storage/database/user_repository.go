package database

import (
	"bombayv/logiapp-monorepo/logi_api/internal/core/user"
	"context"
)

// UserRepository handles the database operations for users.
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// Save saves a user and their associated data to the database in a transaction.
func (r *UserRepository) Save(ctx context.Context, u *user.User, ud *user.UserData) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	userQuery := `
		INSERT INTO users (user_id, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, userQuery, u.UserID, u.Email, u.PasswordHash, u.Role, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}

	userDataQuery := `
		INSERT INTO users_data (user_id, first_name, last_name, phone_number, last_connection, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, userDataQuery, ud.UserID, ud.FirstName, ud.LastName, ud.PhoneNumber, ud.LastConnection, ud.CreatedAt, ud.UpdatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// FindByEmail retrieves a user by their email address.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT user_id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	u := &user.User{}
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&u.UserID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
