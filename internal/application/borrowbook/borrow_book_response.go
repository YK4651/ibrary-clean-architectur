package borrowbook

import (
	"time"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/loandm"
)

type BorrowBookResponse struct {
	LoanID     string
	UserID     string
	BookID     string
	BookTitle  string
	BorrowedAt time.Time
	DueDate    time.Time
}

func NewBorrowBookResponse(
	loan *loandm.Loan,
	bookTitle string,
) *BorrowBookResponse {
	loanID := loan.Id()
	userID := loan.UserID()
	bookID := loan.BookID()

	return &BorrowBookResponse{
		LoanID:     loanID.Value(),
		UserID:     userID.Value(),
		BookID:     bookID.Value(),
		BookTitle:  bookTitle,
		BorrowedAt: loan.BorrowedAt(),
		DueDate:    loan.DueDate(),
	}
}
