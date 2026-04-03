package bookdm

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=book_repository.go -destination=../../application/borrowbook/mocks/mock_book_repository.go -package=mocks

import "context"

//go:generate mockgen -source=book_repository.go -destination=../../application/borrowbook/mocks/mock_book_repository.go -package=mocks

type BookRepository interface {
	FindByID(ctx context.Context, id *BookID) (*Book, error)
}
