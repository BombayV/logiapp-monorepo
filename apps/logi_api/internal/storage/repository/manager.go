package repository

import (
	"bombayv/logiapp-monorepo/logi_api/internal/core/orders"
	"bombayv/logiapp-monorepo/logi_api/internal/core/user"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/database"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// Manager provides a unified interface to all repositories
// and handles cross-domain operations and transactions
type Manager struct {
	db     *database.DB
	helper *database.DatabaseHelper
	User   user.Repository
	Orders orders.Repository
}

// NewManager creates a new repository manager
func NewManager(db *database.DB) *Manager {
	return &Manager{
		db:     db,
		helper: database.NewDatabaseHelper(db.Pool),
		User:   database.NewUserRepository(db),
		Orders: database.NewOrderRepository(db),
	}
}

// WithTransaction executes a function within a database transaction
// This allows multiple repositories to participate in the same transaction
func (m *Manager) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := m.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// CreateOrderWithUser is an example of a cross-domain operation
// that benefits from transaction management
func (m *Manager) CreateOrderWithUser(ctx context.Context, email, address string) (*orders.Order, error) {
	var order *orders.Order

	err := m.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// Find user by email
		userObj, err := m.findUserByEmailTx(ctx, tx, email)
		if err != nil {
			return err
		}

		// Create order
		order, err = m.createOrderTx(ctx, tx, userObj.UserID, address)
		if err != nil {
			return err
		}

		return nil
	})

	return order, err
}

// Helper methods that work with transactions
func (m *Manager) findUserByEmailTx(ctx context.Context, tx pgx.Tx, email string) (*user.User, error) {
	query := `
		SELECT user_id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	u := &user.User{}
	err := tx.QueryRow(ctx, query, email).Scan(&u.UserID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *Manager) createOrderTx(ctx context.Context, tx pgx.Tx, userID, address string) (*orders.Order, error) {
	orderID := uuid.New().String()
	now := time.Now()

	order := &orders.Order{
		OrderID:         orderID,
		CreatedBy:       userID,
		AssignedTo:      nil,
		DeliveryAddress: address,
		Status:          "pending",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	query := `
		INSERT INTO orders (order_id, created_by, assigned_to, delivery_address, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(ctx, query, order.OrderID, order.CreatedBy, order.AssignedTo, order.DeliveryAddress, order.Status, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetDatabaseHelper returns the database helper for common operations
func (m *Manager) GetDatabaseHelper() *database.DatabaseHelper {
	return m.helper
}
