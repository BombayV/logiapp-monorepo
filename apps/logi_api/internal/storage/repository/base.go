package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

var (
	ErrTransactionAlreadyStarted = errors.New("transaction already started")
	ErrNoActiveTransaction       = errors.New("no active transaction")
)

// Executor represents something that can execute queries (DB or Tx)
// This interface matches the common methods between pgx.Conn and pgx.Tx
type Executor interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// Transactor can begin transactions
type Transactor interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// Database combines both execution and transaction capabilities
type Database interface {
	Executor
	Transactor
}

// UnitOfWork provides transaction management across multiple repositories
type UnitOfWork struct {
	db Database
	tx pgx.Tx
}

// NewUnitOfWork creates a new unit of work
func NewUnitOfWork(db Database) *UnitOfWork {
	return &UnitOfWork{db: db}
}

// Begin starts a transaction
func (uow *UnitOfWork) Begin(ctx context.Context) error {
	if uow.tx != nil {
		return ErrTransactionAlreadyStarted
	}

	tx, err := uow.db.Begin(ctx)
	if err != nil {
		return err
	}

	uow.tx = tx
	return nil
}

// Commit commits the transaction
func (uow *UnitOfWork) Commit(ctx context.Context) error {
	if uow.tx == nil {
		return ErrNoActiveTransaction
	}

	err := uow.tx.Commit(ctx)
	uow.tx = nil
	return err
}

// Rollback rolls back the transaction
func (uow *UnitOfWork) Rollback(ctx context.Context) error {
	if uow.tx == nil {
		return ErrNoActiveTransaction
	}

	err := uow.tx.Rollback(ctx)
	uow.tx = nil
	return err
}

// GetExecutor returns the current executor (transaction if active, otherwise database)
func (uow *UnitOfWork) GetExecutor() Executor {
	if uow.tx != nil {
		return uow.tx
	}
	return uow.db
}

// WithTransaction executes a function within a transaction
func (uow *UnitOfWork) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := uow.Begin(ctx); err != nil {
		return err
	}

	defer func() {
		if uow.tx != nil {
			uow.Rollback(ctx)
		}
	}()

	if err := fn(ctx); err != nil {
		return err
	}

	return uow.Commit(ctx)
}

// BaseRepository provides common repository functionality
type BaseRepository struct {
	uow *UnitOfWork
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(uow *UnitOfWork) *BaseRepository {
	return &BaseRepository{uow: uow}
}

// GetExecutor returns the current executor
func (r *BaseRepository) GetExecutor() Executor {
	return r.uow.GetExecutor()
}
