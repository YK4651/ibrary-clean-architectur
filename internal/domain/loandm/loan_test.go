package loandm_test

import (
	"testing"
	"time"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/loandm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

func TestDueDateIsCalculatedCorrectly(t *testing.T) {
	borrowedAt := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Create Value Objects
	loanID := loandm.NewLoanID()
	userID := userdm.NewUserID()
	bookID := bookdm.NewBookID()

	l := loandm.NewLoan(
		loanID,
		userID,
		bookID,
		borrowedAt,
		nil, // not returned yet
	)

	expectedDueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	if !l.DueDate().Equal(expectedDueDate) {
		t.Errorf("Expected due date %v, got %v", expectedDueDate, l.DueDate())
	}
}

func TestLoanIsOverdueAfterDueDate(t *testing.T) {
	borrowedAt := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	currentDate := time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC)

	// Create Value Objects
	loanID := loandm.NewLoanID()
	userID := userdm.NewUserID()
	bookID := bookdm.NewBookID()

	l := loandm.NewLoan(loanID, userID, bookID, borrowedAt, nil)

	if !l.IsOverdue(&currentDate) {
		t.Error("Expected loan to be overdue")
	}
}

func TestCanMarkLoanAsReturned(t *testing.T) {
	// Create Value Objects
	loanID := loandm.NewLoanID()
	userID := userdm.NewUserID()
	bookID := bookdm.NewBookID()

	l := loandm.NewLoan(loanID, userID, bookID, time.Now(), nil)

	if l.IsReturned() {
		t.Error("Loan should not be returned initially")
	}

	returned, err := l.MarkAsReturned(nil) // Pass nil to use current time
	if err != nil {
		t.Fatalf("Failed to mark as returned: %v", err)
	}

	if !returned.IsReturned() {
		t.Error("Loan should be marked as returned")
	}
	if returned.ReturnedAt() == nil {
		t.Error("ReturnedAt should be set")
	}
}
