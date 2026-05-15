package models

import (
	"strings"
	"testing"
)

func TestProjectValidatePositive(t *testing.T) {
	project := Project{Name: "Study", Description: "Course project", OwnerID: 1}
	if err := project.Validate(); err != nil {
		t.Fatalf("expected valid project, got error: %v", err)
	}
}

func TestProjectValidateNegativeOwner(t *testing.T) {
	project := Project{Name: "Study", OwnerID: 0}
	if err := project.Validate(); err == nil {
		t.Fatal("expected owner_id validation error")
	}
}

func TestProjectValidateBoundaryName(t *testing.T) {
	tests := []string{
		strings.Repeat("a", 3),
		strings.Repeat("a", 80),
	}
	for _, name := range tests {
		project := Project{Name: name, OwnerID: 1}
		if err := project.Validate(); err != nil {
			t.Fatalf("expected boundary name %q to be valid: %v", name, err)
		}
	}
}
