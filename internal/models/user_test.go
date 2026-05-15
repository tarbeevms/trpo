package models

import (
	"strings"
	"testing"
)

func TestUserValidatePositive(t *testing.T) {
	user := User{Login: "ivan"}
	if err := user.Validate(); err != nil {
		t.Fatalf("expected valid user, got error: %v", err)
	}
}

func TestUserValidateNegativeLoginWithSpace(t *testing.T) {
	user := User{Login: "bad login"}
	if err := user.Validate(); err == nil {
		t.Fatal("expected login validation error")
	}
}

func TestUserValidateBoundaryLogin(t *testing.T) {
	tests := []string{
		strings.Repeat("a", 3),
		strings.Repeat("a", 50),
	}
	for _, login := range tests {
		user := User{Login: login}
		if err := user.Validate(); err != nil {
			t.Fatalf("expected boundary login %q to be valid: %v", login, err)
		}
	}
}
