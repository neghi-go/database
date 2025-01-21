package database

import (
	"context"
)

type Model[T any] interface {
	WithContext(ctx context.Context) Model[T]
	Query(query_params ...Params) Query[T]
	Save(doc ...T) error
	ExecRaw() error
}

// Query interface defines the structure of the store queries
type Query[T any] interface {
	// Count returns the number of documents that match a query
	Count() (int64, error)
	// First returns the first document that matches a query
	First() (*T, error)
	// All returns all the document that matches a query
	All() ([]*T, error)
	// Update updates the document that matches a query
	Update(doc T) error
	// UpdateMany updates all the document that matches a query
	UpdateMany(doc T) error
	// Delete deletes the document that matches a query
	Delete() error
	// DeleteMany deletes all document that matches the query
	DeleteMany() error
}
