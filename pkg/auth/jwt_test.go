package auth

import (
	"testing"
	"time"
)

func TestManagerGenerateAndValidate(t *testing.T) {
	manager := NewManager("secret", time.Hour)
	token, err := manager.Generate(10, "ivan")
	if err != nil {
		t.Fatalf("expected token, got error: %v", err)
	}
	claims, err := manager.Validate(token)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}
	if claims.UserID != 10 || claims.Login != "ivan" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}

func TestManagerValidateNegativeSignature(t *testing.T) {
	manager := NewManager("secret", time.Hour)
	token, err := manager.Generate(10, "ivan")
	if err != nil {
		t.Fatalf("expected token, got error: %v", err)
	}
	otherManager := NewManager("other-secret", time.Hour)
	if _, err := otherManager.Validate(token); err == nil {
		t.Fatal("expected signature validation error")
	}
}
