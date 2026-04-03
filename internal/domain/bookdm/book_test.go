package bookdm_test

import (
	"testing"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
)

func TestCanCreateBook(t *testing.T) {
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}

	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	if len(b.Id().Value()) != 26 {
		t.Errorf("Expected ULID length 26, got %d", len(b.Id().Value()))
	}
	if b.Title() != "Clean Architecture" {
		t.Errorf("Expected title 'Clean Architecture', got '%s'", b.Title())
	}
	if b.TotalCopies() != 3 {
		t.Errorf("Expected 3 total copies, got %d", b.TotalCopies())
	}
}

func TestCanBorrowWhenCopiesAvailable(t *testing.T) {
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 2)
	if err != nil {
		t.Fatal(err)
	}

	// アクティブ貸出0件なら在庫あり（totalCopies=2）
	if !b.IsAvailable(0) {
		t.Error("Expected book to be available when no active loans")
	}
}

func TestIsAvailableWithActiveLoans(t *testing.T) {
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 2)
	if err != nil {
		t.Fatal(err)
	}

	// アクティブ貸出1件 → まだ借りられる
	if !b.IsAvailable(1) {
		t.Error("Expected available when 1 active loan and totalCopies=2")
	}
	// アクティブ貸出2件以上 → 借りられない
	if b.IsAvailable(2) {
		t.Error("Expected not available when active loans >= totalCopies")
	}
}

func TestIsAvailableWhenNoCopiesLeft(t *testing.T) {
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 1)
	if err != nil {
		t.Fatal(err)
	}

	// 1冊のみで1件貸出中 → 借りられない
	if b.IsAvailable(1) {
		t.Error("Expected not available when active loans >= totalCopies")
	}
}
