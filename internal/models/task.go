package models

import (
	"fmt"
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
)

type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

type Task struct {
	BaseEntity
	AuditInfo
	SoftDelete
	ProjectID   int64
	Title       string
	Description string
	Status      Status
	Priority    Priority
	Deadline    time.Time
	AssigneeID  int64
	Tags        []Tag
}

type TaskFilter struct {
	AssigneeID int64
}

type Tag struct {
	ID   int64
	Name string
}

type TaskHistory struct {
	BaseEntity
	TaskID    int64
	OldStatus Status
	NewStatus Status
	ChangedAt time.Time
}

func (t Task) Validate(now time.Time) error {
	if t.ProjectID <= 0 {
		return fmt.Errorf("project_id is required")
	}
	if err := validateLength("title", t.Title, 3, 100); err != nil {
		return err
	}
	if err := validateMaxLength("description", t.Description, 1000); err != nil {
		return err
	}
	if !t.Status.IsValid() {
		return fmt.Errorf("status must be one of: new, in_progress, done, cancelled")
	}
	if !t.Priority.IsValid() {
		return fmt.Errorf("priority must be one of: low, medium, high, critical")
	}
	if dateOnly(t.Deadline).Before(dateOnly(now)) {
		return fmt.Errorf("deadline cannot be in the past")
	}
	if t.AssigneeID <= 0 {
		return fmt.Errorf("assignee_id is required")
	}
	if len(t.Tags) > 10 {
		return fmt.Errorf("task can contain no more than 10 tags")
	}
	for _, tag := range t.Tags {
		if err := tag.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) ChangeStatus(newStatus Status, changedAt time.Time) (TaskHistory, error) {
	if !newStatus.IsValid() {
		return TaskHistory{}, fmt.Errorf("status must be one of: new, in_progress, done, cancelled")
	}
	history := TaskHistory{
		TaskID:    t.ID,
		OldStatus: t.Status,
		NewStatus: newStatus,
		ChangedAt: changedAt,
	}
	t.Status = newStatus
	return history, nil
}

// GetTaskInfo возвращает информацию о задаче, используя методы встроенных структур
// Демонстрирует множественное наследование и разрешение конфликтов методов
func (t Task) GetTaskInfo() string {
	// Использование уникальных методов - работают напрямую
	creatorID := t.GetCreatedBy() // метод из AuditInfo
	isDeleted := t.IsDeleted()    // метод из SoftDelete
	createdAt := t.GetCreatedAt() // метод из BaseEntity

	// Использование методов с одинаковым названием - требует явного указания
	entityID := t.BaseEntity.GetID()          // ID задачи из BaseEntity
	creatorIDFromGetID := t.AuditInfo.GetID() // ID создателя из AuditInfo
	deleteID := t.SoftDelete.GetID()          // статус удаления из SoftDelete

	return fmt.Sprintf("Task[id=%d, creator=%d, deleted=%v, created=%v, entityID=%d, creatorID=%d, deleteStatus=%d]",
		t.ID, creatorID, isDeleted, createdAt, entityID, creatorIDFromGetID, deleteID)
}

func (s Status) IsValid() bool {
	return s == StatusNew || s == StatusInProgress || s == StatusDone || s == StatusCancelled
}

func (p Priority) IsValid() bool {
	return p == PriorityLow || p == PriorityMedium || p == PriorityHigh || p == PriorityCritical
}

func (t Tag) Validate() error {
	return validateLength("tag name", t.Name, 2, 30)
}
