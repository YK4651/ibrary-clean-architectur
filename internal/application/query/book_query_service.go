package query

import "context"

//go:generate mockgen -source=book_query_service.go -destination=mocks/mock_book_query_service.go -package=mocks

type BookQueryService interface {
	// GetBookByID retrieves a book by ID with current loan status
	GetBookByID(ctx context.Context, bookID string) (*BookReadModel, error)

	// ListBooks retrieves all books with their loan status (for Lesson 6)
	ListBooks(ctx context.Context) ([]*BookReadModel, error)
}
