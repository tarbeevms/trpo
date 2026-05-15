package models

import (
	"strings"
	"testing"
	"time"
)

func validTask(now time.Time) Task {
	return Task{
		ProjectID:   1,
		Title:       "Write report",
		Description: "Prepare project description",
		Status:      StatusNew,
		Priority:    PriorityHigh,
		Deadline:    now,
		AssigneeID:  1,
		Tags:        []Tag{{Name: "study"}},
	}
}

func TestTaskValidatePositive(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	task := validTask(now)
	if err := task.Validate(now); err != nil {
		t.Fatalf("expected valid task, got error: %v", err)
	}
}

func TestTaskValidateNegativeStatus(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	task := validTask(now)
	task.Status = "bad"
	if err := task.Validate(now); err == nil {
		t.Fatal("expected status validation error")
	}
}

func TestTaskValidateNegativeDeadline(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	task := validTask(now)
	task.Deadline = now.AddDate(0, 0, -1)
	if err := task.Validate(now); err == nil {
		t.Fatal("expected deadline validation error")
	}
}

func TestTaskValidateBoundaryFields(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	task := validTask(now)
	task.Title = strings.Repeat("a", 3)
	task.Description = strings.Repeat("b", 1000)
	task.Tags = make([]Tag, 10)
	for i := range task.Tags {
		task.Tags[i] = Tag{Name: "tag" + string(rune('a'+i))}
	}
	if err := task.Validate(now); err != nil {
		t.Fatalf("expected minimum title, maximum description, 10 tags and today deadline to be valid: %v", err)
	}

	task.Title = strings.Repeat("a", 100)
	if err := task.Validate(now); err != nil {
		t.Fatalf("expected maximum title to be valid: %v", err)
	}
}

func TestTaskValidateNegativeTooManyTags(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	task := validTask(now)
	task.Tags = make([]Tag, 11)
	for i := range task.Tags {
		task.Tags[i] = Tag{Name: "tag" + string(rune('a'+i))}
	}
	if err := task.Validate(now); err == nil {
		t.Fatal("expected too many tags validation error")
	}
}

func TestTaskChangeStatus(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	task := validTask(now)
	task.ID = 10

	history, err := task.ChangeStatus(StatusDone, now)
	if err != nil {
		t.Fatalf("expected status change, got error: %v", err)
	}
	if task.Status != StatusDone {
		t.Fatalf("expected status %q, got %q", StatusDone, task.Status)
	}
	if history.OldStatus != StatusNew || history.NewStatus != StatusDone {
		t.Fatalf("unexpected history: %#v", history)
	}
}

func TestTagValidateBoundary(t *testing.T) {
	for _, name := range []string{strings.Repeat("a", 2), strings.Repeat("a", 30)} {
		if err := (Tag{Name: name}).Validate(); err != nil {
			t.Fatalf("expected boundary tag %q to be valid: %v", name, err)
		}
	}
}
