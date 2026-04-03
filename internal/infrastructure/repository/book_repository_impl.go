package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
)

type BookRepositoryImpl struct{}

func NewBookRepositoryImpl() bookdm.BookRepository {
	return &BookRepositoryImpl{}
}

func (r *BookRepositoryImpl) FindByID(ctx context.Context, id *bookdm.BookID) (*bookdm.Book, error) {
	db, err := getDBExecutor(ctx)
	if err != nil {
		return nil, err
	}

	query := "SELECT id, title, author, isbn FROM books WHERE id = ?"

	var idStr, title, author, isbnStr string
	err = db.QueryRowContext(ctx, query, id.Value()).Scan(
		&idStr, &title, &author, &isbnStr,
	)
	_ = idStr // idは引数から既知
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query book: %w", err)
	}

	isbn, err := bookdm.NewISBN(isbnStr)
	if err != nil {
		return nil, err
	}
	// total_copiesはスキーマにないため1を使用
	return bookdm.NewBook(*id, title, author, isbn, 1)
}
