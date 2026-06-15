package bookdm_test

import (
	"testing"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
)

func TestCanCreateBook(t *testing.T) {
	// Create Value Objects
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}

	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	// ID should be auto-generated ULID (26 characters)
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

func TestIsAvailableWhenCopiesAvailable(t *testing.T) {
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 2)
	if err != nil {
		t.Fatal(err)
	}

	if !b.IsAvailable(0) {
		t.Error("Expected book to be available when no active loans")
	}
	if !b.IsAvailable(1) {
		t.Error("Expected book to be available when active loans < total copies")
	}
}

func TestNotAvailableWhenAllCopiesLoaned(t *testing.T) {
	bookID := bookdm.NewBookID()
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}
	b, err := bookdm.NewBook(bookID, "Clean Architecture", "Robert Martin", isbn, 2)
	if err != nil {
		t.Fatal(err)
	}

	if b.IsAvailable(2) {
		t.Error("Expected book to be unavailable when active loans == total copies")
	}
	if b.IsAvailable(3) {
		t.Error("Expected book to be unavailable when active loans > total copies")
	}
}