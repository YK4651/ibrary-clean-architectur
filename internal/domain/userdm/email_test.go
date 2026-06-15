package userdm_test

import (
	"testing"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

func TestEmail_ValidEmail(t *testing.T) {
	email, err := userdm.NewEmail("john@example.com")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if email.Value() != "john@example.com" {
		t.Errorf("Expected john@example.com, got: %s", email.Value())
	}
}

func TestEmail_InvalidFormat(t *testing.T) {
	_, err := userdm.NewEmail("invalid")
	if err == nil {
		t.Fatal("Expected error for invalid email")
	}
}

func TestEmail_GetDomain(t *testing.T) {
	email, err := userdm.NewEmail("user@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if email.Domain() != "example.com" {
		t.Errorf("Expected example.com, got: %s", email.Domain())
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, err := userdm.NewEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	email2, err := userdm.NewEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	email3, err := userdm.NewEmail("other@example.com")
	if err != nil {
		t.Fatal(err)
	}
	
	if !email1.Equals(email2) {
		t.Error("Expected equal emails")
	}
	if email1.Equals(email3) {
		t.Error("Expected different emails")
	}
}