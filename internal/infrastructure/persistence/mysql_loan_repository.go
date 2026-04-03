package persistence

import (
	"context"
	"database/sql"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/loandm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

type LoanRepositoryImpl struct {
	db *sql.DB
}

func NewLoanRepositoryImpl(db *sql.DB) loandm.LoanRepository {
	return &LoanRepositoryImpl{db: db}
}

func (r *LoanRepositoryImpl) CountActiveLoansForUser(ctx context.Context, userID *userdm.UserID) (uint32, error) {
	query := "SELECT COUNT(*) FROM loans WHERE user_id = ? AND returned_at IS NULL"

	var count uint32
	err := r.db.QueryRowContext(ctx, query, userID.Value()).Scan(&count)
	return count, err
}

func (r *LoanRepositoryImpl) CountActiveLoansForBook(ctx context.Context, bookID *bookdm.BookID) (uint32, error) {
	query := "SELECT COUNT(*) FROM loans WHERE book_id = ? AND returned_at IS NULL"

	var count uint32
	err := r.db.QueryRowContext(ctx, query, bookID.Value()).Scan(&count)
	return count, err
}

func (r *LoanRepositoryImpl) Save(ctx context.Context, loan *loandm.Loan) error {
	query := `INSERT INTO loans (id, user_id, book_id, borrowed_at, due_date, returned_at)
	          VALUES (?, ?, ?, ?, ?, ?)`

	var returnedAt *string
	if loan.ReturnedAt() != nil {
		t := loan.ReturnedAt().Format("2006-01-02 15:04:05")
		returnedAt = &t
	}

	loanID := loan.Id()
	userID := loan.UserID()
	bookID := loan.BookID()

	_, err := r.db.ExecContext(ctx, query,
		loanID.Value(),
		userID.Value(),
		bookID.Value(),
		loan.BorrowedAt().Format("2006-01-02 15:04:05"),
		loan.DueDate().Format("2006-01-02 15:04:05"),
		returnedAt,
	)
	return err
}
