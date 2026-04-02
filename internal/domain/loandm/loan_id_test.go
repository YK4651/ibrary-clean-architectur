package loandm_test

import (
	"testing"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/loandm"
)

func TestLoanID_GeneratesValidID(t *testing.T) {
	id := loandm.NewLoanID()

	if id.Value() == "" {
		t.Error("Expected non-empty ID value")
	}

	// ULID should be 26 characters
	if len(id.Value()) != 26 {
		t.Errorf("Expected ULID length 26, got: %d", len(id.Value()))
	}
}

func TestLoanID_GeneratesUniqueIDs(t *testing.T) {
	id1 := loandm.NewLoanID()
	id2 := loandm.NewLoanID()

	if id1.Equals(id2) {
		t.Error("Expected different IDs to not be equal")
	}
}

func TestLoanID_Equality(t *testing.T) {
	id1 := loandm.NewLoanID()
	id2 := id1 // Same reference

	if !id1.Equals(id2) {
		t.Error("Expected same ID to be equal to itself")
	}
}

func TestLoanID_ValueReturnsString(t *testing.T) {
	id := loandm.NewLoanID()
	value := id.Value()

	if value == "" {
		t.Error("Value() should return non-empty string")
	}
}
