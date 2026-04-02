package loandm_test

import (
	"testing"
	"time"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/loandm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

func TestUserCanBorrowWhenAllConditionsMet(t *testing.T) {
	service := loandm.NewLoanEligibilityService()

	// Create user with default state (0 loans, 0 fees)
	u := userdm.NewUser("John Doe", "john@example.com")

	// Create book
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatal(err)
	}

	// currentLoanCount=0, currentActiveLoans=0
	if !service.CanBorrow(u, 0, b, 0) {
		t.Error("ユーザーは本を借りられるべきです")
	}
}

func TestUserCannotBorrowWhenMaxLoansReached(t *testing.T) {
	service := loandm.NewLoanEligibilityService()

	// Create user using ReconstructUser
	userID := userdm.NewUserID()
	u := userdm.ReconstructUser(
		userID,
		"John Doe",
		"john@example.com",
		userdm.UserStatusActive,
		0, // overdueFees
		time.Now(),
	)

	// Create book
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatal(err)
	}

	// currentLoanCount=MaxLoans (user at limit)
	if service.CanBorrow(u, userdm.MaxLoans, b, 0) {
		t.Error("ユーザーは本を借りられないべきです（貸出上限に達している）")
	}
}
