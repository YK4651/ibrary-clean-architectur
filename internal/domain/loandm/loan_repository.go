package loandm

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=loan_repository.go -destination=../../application/borrowbook/mocks/mock_loan_repository.go -package=mocks

import (
	"context"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

//go:generate mockgen -source=loan_repository.go -destination=../../application/borrowbook/mocks/mock_loan_repository.go -package=mocks

type LoanRepository interface {
	// Loanテーブルからアクティブな貸出をカウント
	CountActiveLoansForUser(ctx context.Context, userID *userdm.UserID) (uint32, error)
	CountActiveLoansForBook(ctx context.Context, bookID *bookdm.BookID) (uint32, error)

	Save(ctx context.Context, loan *Loan) error
}
