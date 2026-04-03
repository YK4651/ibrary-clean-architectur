package bookdm

import "context"

type BookRepository interface {
	FindByID(ctx context.Context, id *BookID) (*Book, error)
}
