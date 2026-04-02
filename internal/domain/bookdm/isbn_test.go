package bookdm_test

import (
	"testing"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
)

func TestCanCreateValidISBN(t *testing.T) {
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatalf("Failed to create ISBN: %v", err)
	}

	if isbn.Value() != "9780134494166" {
		t.Errorf("Expected '9780134494166', got '%s'", isbn.Value())
	}
}

func TestCanCreateISBNWithHyphens(t *testing.T) {
	isbn, err := bookdm.NewISBN("978-0-13-449416-6")
	if err != nil {
		t.Fatalf("Failed to create ISBN: %v", err)
	}

	if isbn.Value() != "9780134494166" {
		t.Errorf("Expected normalized '9780134494166', got '%s'", isbn.Value())
	}
}

func TestInvalidISBNReturnsError(t *testing.T) {
	_, err := bookdm.NewISBN("invalid")
	if err == nil {
		t.Error("Expected error for invalid ISBN")
	}
}

func TestFormatsISBNCorrectly(t *testing.T) {
	isbn, err := bookdm.NewISBN("9780134494166")
	if err != nil {
		t.Fatal(err)
	}

	if isbn.Formatted() != "978-0-13-449416-6" {
		t.Errorf("Expected '978-0-13-449416-6', got '%s'", isbn.Formatted())
	}
}
