package borrowbook_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/YK4651/ibrary-clean-architectur/internal/application/borrowbook"
	"github.com/YK4651/ibrary-clean-architectur/internal/application/borrowbook/mocks"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/loandm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
	"go.uber.org/mock/gomock"
)

func TestBorrowBookUseCase_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockBookRepo := mocks.NewMockBookRepository(ctrl)
	mockLoanRepo := mocks.NewMockLoanRepository(ctrl)

	ctx := context.Background()
	userID, err := userdm.NewUserID("12345678")
	if err != nil {
		t.Fatal(err)
	}
	u := userdm.ReconstructUser(userID, "John", "john@example.com", userdm.UserStatusActive, 0, time.Now())

	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatal(err)
	}

	// モックの期待値を設定
	mockUserRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(u, nil)
	mockBookRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(b, nil)
	mockLoanRepo.EXPECT().CountActiveLoansForUser(ctx, gomock.Any()).Return(uint32(0), nil)
	mockLoanRepo.EXPECT().CountActiveLoansForBook(ctx, gomock.Any()).Return(uint32(0), nil)
	mockLoanRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)

	useCase := borrowbook.NewBorrowBookUseCase(mockUserRepo, mockBookRepo, mockLoanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookID.Value())
	if err != nil {
		t.Fatal(err)
	}

	// Act
	resp, err := useCase.Execute(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	if resp.LoanID == "" {
		t.Error("Expected loan ID, got empty string")
	}
}

func TestBorrowBookUseCase_UserNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockBookRepo := mocks.NewMockBookRepository(ctrl)
	mockLoanRepo := mocks.NewMockLoanRepository(ctrl)

	ctx := context.Background()
	bookID := bookdm.NewBookID()

	mockUserRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, nil)

	useCase := borrowbook.NewBorrowBookUseCase(mockUserRepo, mockBookRepo, mockLoanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookID.Value())
	if err != nil {
		t.Fatal(err)
	}

	// Act
	_, err = useCase.Execute(ctx, req)

	// Assert
	if err == nil {
		t.Error("Expected error when user not found")
	}
	var target *borrowbook.UserNotFoundError
	if !errors.As(err, &target) {
		t.Errorf("Expected UserNotFoundError, got %T", err)
	}
}

func TestBorrowBookUseCase_BookNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockBookRepo := mocks.NewMockBookRepository(ctrl)
	mockLoanRepo := mocks.NewMockLoanRepository(ctrl)

	ctx := context.Background()
	userID, err := userdm.NewUserID("12345678")
	if err != nil {
		t.Fatal(err)
	}
	u := userdm.ReconstructUser(userID, "John", "john@example.com", userdm.UserStatusActive, 0, time.Now())
	bookID := bookdm.NewBookID()

	mockUserRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(u, nil)
	mockBookRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, nil)

	useCase := borrowbook.NewBorrowBookUseCase(mockUserRepo, mockBookRepo, mockLoanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookID.Value())
	if err != nil {
		t.Fatal(err)
	}

	// Act
	_, err = useCase.Execute(ctx, req)

	// Assert
	if err == nil {
		t.Error("Expected error when book not found")
	}
	var target *borrowbook.BookNotFoundError
	if !errors.As(err, &target) {
		t.Errorf("Expected BookNotFoundError, got %T", err)
	}
}

func TestBorrowBookUseCase_LoanLimitExceeded(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockBookRepo := mocks.NewMockBookRepository(ctrl)
	mockLoanRepo := mocks.NewMockLoanRepository(ctrl)

	ctx := context.Background()
	userID, err := userdm.NewUserID("12345678")
	if err != nil {
		t.Fatal(err)
	}
	u := userdm.ReconstructUser(userID, "John", "john@example.com", userdm.UserStatusActive, 0, time.Now())

	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatal(err)
	}

	mockUserRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(u, nil)
	mockBookRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(b, nil)
	mockLoanRepo.EXPECT().CountActiveLoansForUser(ctx, gomock.Any()).Return(uint32(5), nil) // 上限に達している
	mockLoanRepo.EXPECT().CountActiveLoansForBook(ctx, gomock.Any()).Return(uint32(0), nil)

	useCase := borrowbook.NewBorrowBookUseCase(mockUserRepo, mockBookRepo, mockLoanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookID.Value())
	if err != nil {
		t.Fatal(err)
	}

	// Act
	_, err = useCase.Execute(ctx, req)

	// Assert
	if err == nil {
		t.Error("Expected UserCannotBorrowError when loan limit reached")
	}
	var target *borrowbook.UserCannotBorrowError
	if !errors.As(err, &target) {
		t.Errorf("Expected UserCannotBorrowError, got %T", err)
	}
}

func TestBorrowBookUseCase_BookNotAvailable(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockBookRepo := mocks.NewMockBookRepository(ctrl)
	mockLoanRepo := mocks.NewMockLoanRepository(ctrl)

	ctx := context.Background()
	userID, err := userdm.NewUserID("12345678")
	if err != nil {
		t.Fatal(err)
	}
	u := userdm.ReconstructUser(userID, "John", "john@example.com", userdm.UserStatusActive, 0, time.Now())

	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 1)
	if err != nil {
		t.Fatal(err)
	}

	mockUserRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(u, nil)
	mockBookRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(b, nil)
	mockLoanRepo.EXPECT().CountActiveLoansForUser(ctx, gomock.Any()).Return(uint32(0), nil)
	mockLoanRepo.EXPECT().CountActiveLoansForBook(ctx, gomock.Any()).Return(uint32(1), nil) // 全コピー貸出中

	useCase := borrowbook.NewBorrowBookUseCase(mockUserRepo, mockBookRepo, mockLoanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookID.Value())
	if err != nil {
		t.Fatal(err)
	}

	// Act
	_, err = useCase.Execute(ctx, req)

	// Assert
	if err == nil {
		t.Error("Expected BookNotAvailableError")
	}
	var target *borrowbook.BookNotAvailableError
	if !errors.As(err, &target) {
		t.Errorf("Expected BookNotAvailableError, got %T", err)
	}
}

func TestBorrowBookUseCase_CreatesLoanWithCorrectDueDate(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockBookRepo := mocks.NewMockBookRepository(ctrl)
	mockLoanRepo := mocks.NewMockLoanRepository(ctrl)

	ctx := context.Background()
	userID, err := userdm.NewUserID("12345678")
	if err != nil {
		t.Fatal(err)
	}
	u := userdm.ReconstructUser(userID, "John", "john@example.com", userdm.UserStatusActive, 0, time.Now())

	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatal(err)
	}

	mockUserRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(u, nil)
	mockBookRepo.EXPECT().FindByID(ctx, gomock.Any()).Return(b, nil)
	mockLoanRepo.EXPECT().CountActiveLoansForUser(ctx, gomock.Any()).Return(uint32(0), nil)
	mockLoanRepo.EXPECT().CountActiveLoansForBook(ctx, gomock.Any()).Return(uint32(0), nil)

	var savedLoan *loandm.Loan
	mockLoanRepo.EXPECT().Save(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, loan *loandm.Loan) error {
			savedLoan = loan
			return nil
		},
	)

	useCase := borrowbook.NewBorrowBookUseCase(mockUserRepo, mockBookRepo, mockLoanRepo)

	req, err := borrowbook.NewBorrowBookRequest("12345678", bookID.Value())
	if err != nil {
		t.Fatal(err)
	}

	// Act
	_, err = useCase.Execute(ctx, req)

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if savedLoan == nil {
		t.Fatal("Expected loan to be saved")
	}
	dueDate := savedLoan.DueDate()
	now := time.Now()
	expectedDue := now.AddDate(0, 0, loandm.LoanPeriodDays)
	if dueDate.Format("2006-01-02") != expectedDue.Format("2006-01-02") {
		t.Errorf("Expected due date %v, got %v", expectedDue.Format("2006-01-02"), dueDate.Format("2006-01-02"))
	}
}
