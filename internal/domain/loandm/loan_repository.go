package loandm

import (
	"context"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

type LoanRepository interface {
	// Loanテーブルからアクティブな貸出をカウント
	CountActiveLoansForUser(ctx context.Context, userID *userdm.UserID) (uint32, error)
	CountActiveLoansForBook(ctx context.Context, bookID *bookdm.BookID) (uint32, error)

	Save(ctx context.Context, loan *Loan) error
}
