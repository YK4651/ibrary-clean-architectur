package borrowbook

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=borrow_book_usecase.go -destination=mocks/mock_repositories.go -package=mocks

// ↓ 以下は Step 4 で書いた既存コード（変更不要）
import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/loandm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

type BorrowBookUseCase struct {
	userRepo           UserRepository
	bookRepo           bookdm.BookRepository
	loanRepo           loandm.LoanRepository
	eligibilityService *loandm.LoanEligibilityService
}

func NewBorrowBookUseCase(
	userRepo UserRepository,
	bookRepo bookdm.BookRepository,
	loanRepo loandm.LoanRepository,
) *BorrowBookUseCase {
	return &BorrowBookUseCase{
		userRepo:           userRepo,
		bookRepo:           bookRepo,
		loanRepo:           loanRepo,
		eligibilityService: loandm.NewLoanEligibilityService(),
	}
}

// Execute - 書籍を借りる - 複数エンティティを調整
// ミドルウェアがトランザクションを自動管理
func (uc *BorrowBookUseCase) Execute(ctx context.Context, req *BorrowBookRequest) (*BorrowBookResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	userID, err := userdm.NewUserID(req.UserID)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, &UserNotFoundError{UserID: req.UserID}
	}

	bookID, err := bookdm.BookIDFromString(req.BookID)
	if err != nil {
		return nil, err
	}

	book, err := uc.bookRepo.FindByID(ctx, &bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to find book: %w", err)
	}
	if book == nil {
		return nil, &BookNotFoundError{BookID: req.BookID}
	}

	userCurrentLoans, err := uc.loanRepo.CountActiveLoansForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count user loans: %w", err)
	}

	bookActiveLoans, err := uc.loanRepo.CountActiveLoansForBook(ctx, &bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to count book loans: %w", err)
	}

	if !uc.eligibilityService.CanBorrow(user, userCurrentLoans, book, bookActiveLoans) {
		reason := "user cannot borrow this book"
		if ineligibilityReason := uc.eligibilityService.IneligibilityReason(user, userCurrentLoans, book, bookActiveLoans); ineligibilityReason != nil {
			reason = *ineligibilityReason
		}

		if !user.CanBorrow(userCurrentLoans) {
			return nil, &UserCannotBorrowError{
				UserID: req.UserID,
				Reason: reason,
			}
		}

		return nil, &BookNotAvailableError{
			BookID: req.BookID,
			Reason: reason,
		}
	}

	loan := loandm.NewLoan(
		loandm.NewLoanID(),
		*userID,
		bookID,
		time.Now(),
		nil,
	)

	if err := uc.loanRepo.Save(ctx, loan); err != nil {
		return nil, fmt.Errorf("failed to save loan: %w", err)
	}

	return NewBorrowBookResponse(loan, book.Title()), nil
}

type UserNotFoundError struct {
	UserID string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user not found: %s", e.UserID)
}

type BookNotFoundError struct {
	BookID string
}

func (e *BookNotFoundError) Error() string {
	return fmt.Sprintf("book not found: %s", e.BookID)
}

type UserCannotBorrowError struct {
	UserID string
	Reason string
}

func (e *UserCannotBorrowError) Error() string {
	if e.Reason != "" {
		return e.Reason
	}
	return fmt.Sprintf("user cannot borrow: %s", e.UserID)
}

type BookNotAvailableError struct {
	BookID string
	Reason string
}

func (e *BookNotAvailableError) Error() string {
	if e.Reason != "" {
		return e.Reason
	}
	return fmt.Sprintf("book is not available: %s", e.BookID)
}
