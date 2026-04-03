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

	u := userdm.NewUser("John Doe", "john@example.com")

	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatal(err)
	}

	// ユーザー貸出0件、当該書籍の貸出0件 → 借りられる
	if !service.CanBorrow(u, 0, b, 0) {
		t.Error("ユーザーは本を借りられるべきです")
	}
}

func TestUserCannotBorrowWhenMaxLoansReached(t *testing.T) {
	service := loandm.NewLoanEligibilityService()

	userID := userdm.GenerateUserID()
	u := userdm.ReconstructUser(
		userID,
		"John Doe",
		"john@example.com",
		userdm.UserStatusActive,
		0,
		time.Now(),
	)

	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatal(err)
	}

	// currentLoanCount=5 (MaxLoans) => 借りられない
	if service.CanBorrow(u, 5, b, 0) {
		t.Error("ユーザーは本を借りられないべきです（貸出上限に達している）")
	}
}
