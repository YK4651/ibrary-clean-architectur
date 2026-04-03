package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/YK4651/ibrary-clean-architectur/internal/application/borrowbook"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

type dbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func getDBExecutor(ctx context.Context) (dbExecutor, error) {
	db := ctx.Value("db")
	switch typedDB := db.(type) {
	case *sql.DB:
		return typedDB, nil
	case *sql.Tx:
		return typedDB, nil
	case nil:
		return nil, errors.New("database connection not found in context")
	default:
		return nil, fmt.Errorf("unsupported database connection type: %T", db)
	}
}

type BorrowBookUserRepository struct{}

func NewBorrowBookUserRepository() borrowbook.UserRepository {
	return &BorrowBookUserRepository{}
}

func (r *BorrowBookUserRepository) FindByID(ctx context.Context, id *userdm.UserID) (*userdm.User, error) {
	db, err := getDBExecutor(ctx)
	if err != nil {
		return nil, err
	}

	query := "SELECT id, name, email, status, created_at FROM users WHERE id = ?"

	var idStr, name, email string
	var status uint8
	var createdAt time.Time
	err = db.QueryRowContext(ctx, query, id.Value()).Scan(&idStr, &name, &email, &status, &createdAt)
	_ = idStr // idは引数から既知
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// DBのstatusをドメインのUserStatusにマッピング
	userStatus := userdm.UserStatusActive
	if status == 2 {
		userStatus = userdm.UserStatusSuspended
	}

	return userdm.ReconstructUser(
		id,
		name,
		email,
		userStatus,
		0, // overdue_feesはusersテーブルにないため0
		createdAt,
	), nil
}
