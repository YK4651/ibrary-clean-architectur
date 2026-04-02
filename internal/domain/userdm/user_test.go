package userdm_test

import (
	"testing"
	"time"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

func TestCanCreateValidUser(t *testing.T) {
	// Create user with default state
	u := userdm.NewUser("John Doe", "john@example.com")

	// ID should be auto-generated ULID (26 characters)
	if len(u.Id().Value()) != 26 {
		t.Errorf("Expected UserID length 26, got %d", len(u.Id().Value()))
	}
	if u.Name() != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", u.Name())
	}
	if u.Email() != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", u.Email())
	}
	if u.Status() != userdm.UserStatusActive {
		t.Errorf("Expected status 'active', got '%s'", u.Status())
	}
}

func TestCannotBorrowWhenMaxLoansReached(t *testing.T) {
	// Create user
	userID := userdm.NewUserID()
	u := userdm.ReconstructUser(
		userID,
		"John Doe",
		"john@example.com",
		userdm.UserStatusActive,
		0,
		time.Now(),
	)

	if u.CanBorrow(userdm.MaxLoans) {
		t.Error("ユーザーは貸出上限に達しているため、これ以上借りられないべきです")
	}
}

func TestCanBorrowWhenUnderLimit(t *testing.T) {
	// Create user
	userID := userdm.NewUserID()
	u := userdm.ReconstructUser(
		userID,
		"John Doe",
		"john@example.com",
		userdm.UserStatusActive,
		0,
		time.Now(),
	)

	if !u.CanBorrow(2) {
		t.Error("ユーザーは貸出上限未満のため、さらに借りられるべきです")
	}
}
