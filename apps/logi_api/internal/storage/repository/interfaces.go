package repository

import "context"

// Common repository patterns and interfaces that can be shared

// Searcher provides search capabilities for repositories
type Searcher interface {
	Search(ctx context.Context, query string, limit int) (interface{}, error)
}

// Paginator provides pagination support for repositories
type Paginator interface {
	FindWithPagination(ctx context.Context, limit, offset int) (interface{}, int, error)
}

// SoftDeleter provides soft delete capabilities
type SoftDeleter interface {
	SoftDelete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
}

// TransactionManager provides transaction support
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
