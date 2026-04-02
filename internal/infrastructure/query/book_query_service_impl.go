package query

import (
	"context"
	"database/sql"

	appquery "github.com/YK4651/ibrary-clean-architectur/internal/application/query"
)

// BookQueryServiceImpl - BookQueryServiceインターフェースの実装
//
// 注意：DB接続はフィールドとして保持しません。
// 代わりに、各メソッドでcontextから取得します。
type BookQueryServiceImpl struct {
	// DB フィールドなし - context から取得
}

// NewBookQueryService - 新しいBookQueryServiceImplを作成
//
// DBパラメータは不要 - contextから取得します
func NewBookQueryService() *BookQueryServiceImpl {
	return &BookQueryServiceImpl{}
}

// GetBookByID - IDで書籍を取得し、現在の貸出情報も含める
func (s *BookQueryServiceImpl) GetBookByID(ctx context.Context, bookID string) (*appquery.BookReadModel, error) {
	// ミドルウェアが提供するDB接続をcontextから取得
	db := ctx.Value("db").(*sql.DB)

	// 単一 JOIN クエリで効率的に取得
	q := `
		SELECT
			b.id,
			b.isbn,
			b.title,
			b.author,
			l.id as loan_id,
			l.user_id,
			l.borrowed_at,
			l.due_date
		FROM books b
		LEFT JOIN loans l ON b.id = l.book_id
			AND l.returned_at IS NULL  -- アクティブな貸出のみ
		WHERE b.id = ?
	`

	var (
		id, isbn, title, author string
		loanID, userID          sql.NullString // NULLを許可
		borrowedAt, dueDate     sql.NullString
	)

	// QueryRowContext を使用してcontextを渡す
	err := db.QueryRowContext(ctx, q, bookID).Scan(
		&id, &isbn, &title, &author,
		&loanID, &userID, &borrowedAt, &dueDate,
	)

	if err == sql.ErrNoRows {
		return nil, nil // 書籍が見つかりません
	}
	if err != nil {
		return nil, err // データベースエラー
	}

	// クエリ結果から直接Read Modelを構築
	var currentLoan *appquery.LoanReadModel
	if loanID.Valid { // 貸出がある場合
		currentLoan = &appquery.LoanReadModel{
			LoanID:       loanID.String,
			UserID:       userID.String,
			BorrowedDate: borrowedAt.String,
			DueDate:      dueDate.String,
		}
	}

	return &appquery.BookReadModel{
		ID:          id,
		ISBN:        isbn,
		Title:       title,
		Author:      author,
		IsAvailable: currentLoan == nil, // 貸出データから判断
		CurrentLoan: currentLoan,
	}, nil
}

// ListBooks - すべての書籍を取得（Lesson 5で実装予定）
func (s *BookQueryServiceImpl) ListBooks(ctx context.Context) ([]*appquery.BookReadModel, error) {
	// Lesson 5で実装予定
	return []*appquery.BookReadModel{}, nil
}
