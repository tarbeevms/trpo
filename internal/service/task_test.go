package service

import (
	"context"
	"testing"
	"time"

	"taskflow/internal/models"
)

type fakeLogger struct {
	errors int
}

func (l *fakeLogger) Info(message string, args ...any) {}

func (l *fakeLogger) Error(message string, args ...any) {
	l.errors++
}

type fakeUsers struct {
	exists map[int64]bool
}

func (f fakeUsers) Create(ctx context.Context, user *models.User) error {
	user.ID = 1
	return nil
}

func (f fakeUsers) List(ctx context.Context) ([]models.User, error) {
	return nil, nil
}

func (f fakeUsers) Exists(ctx context.Context, id int64) (bool, error) {
	return f.exists[id], nil
}

func (f fakeUsers) LoginExists(ctx context.Context, login string) (bool, error) {
	return false, nil
}

func (f fakeUsers) FindByID(ctx context.Context, id int64) (models.User, error) {
	return models.User{BaseEntity: models.BaseEntity{ID: id}, Login: "user"}, nil
}

func (f fakeUsers) FindByLogin(ctx context.Context, login string) (models.User, error) {
	return models.User{BaseEntity: models.BaseEntity{ID: 1}, Login: login}, nil
}

type fakeProjects struct {
	exists map[int64]bool
}

func (f fakeProjects) Create(ctx context.Context, project *models.Project) error {
	project.ID = 1
	return nil
}

func (f fakeProjects) List(ctx context.Context) ([]models.Project, error) {
	return nil, nil
}

func (f fakeProjects) ListByOwner(ctx context.Context, ownerID int64) ([]models.Project, error) {
	return nil, nil
}

func (f fakeProjects) Exists(ctx context.Context, id int64) (bool, error) {
	return f.exists[id], nil
}

func (f fakeProjects) OwnedBy(ctx context.Context, projectID int64, ownerID int64) (bool, error) {
	return f.exists[projectID], nil
}

type fakeTasks struct {
	task          models.Task
	created       bool
	statusUpdated bool
}

func (f *fakeTasks) Create(ctx context.Context, task *models.Task) error {
	task.ID = 10
	f.task = *task
	f.created = true
	return nil
}

func (f *fakeTasks) List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
	return []models.Task{f.task}, nil
}

func (f *fakeTasks) FindByID(ctx context.Context, id int64) (models.Task, error) {
	f.task.ID = id
	return f.task, nil
}

func (f *fakeTasks) UpdateStatus(ctx context.Context, task models.Task, history models.TaskHistory) error {
	f.task = task
	f.statusUpdated = true
	return nil
}

func TestTaskFacadeCreatePositive(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	tasks := &fakeTasks{}
	logg := &fakeLogger{}
	facade := NewTaskFacade(
		tasks,
		fakeProjects{exists: map[int64]bool{1: true}},
		fakeUsers{exists: map[int64]bool{2: true}},
		logg,
	)
	facade.SetNow(func() time.Time { return now })

	task := models.Task{
		ProjectID:   1,
		Title:       "Write report",
		Description: "Prepare project description",
		Priority:    models.PriorityHigh,
		Deadline:    now,
		AssigneeID:  2,
		Tags:        []models.Tag{{Name: "study"}},
	}
	created, err := facade.Create(context.Background(), task)
	if err != nil {
		t.Fatalf("expected task creation, got error: %v", err)
	}
	if created.ID != 10 || !tasks.created {
		t.Fatalf("expected repository create to be called, task: %#v", created)
	}
}

func TestTaskFacadeCreateNegativeValidation(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	tasks := &fakeTasks{}
	logg := &fakeLogger{}
	facade := NewTaskFacade(
		tasks,
		fakeProjects{exists: map[int64]bool{1: true}},
		fakeUsers{exists: map[int64]bool{2: true}},
		logg,
	)
	facade.SetNow(func() time.Time { return now })

	task := models.Task{
		ProjectID:  1,
		Title:      "",
		Priority:   models.PriorityHigh,
		Deadline:   now,
		AssigneeID: 2,
	}
	if _, err := facade.Create(context.Background(), task); err == nil {
		t.Fatal("expected validation error")
	}
	if tasks.created {
		t.Fatal("task must not be saved after validation error")
	}
	if logg.errors == 0 {
		t.Fatal("expected error to be logged")
	}
}

func TestTaskFacadeChangeStatusWritesHistory(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	tasks := &fakeTasks{
		task: models.Task{
			BaseEntity: models.BaseEntity{ID: 10},
			ProjectID:  1,
			Title:      "Write report",
			Status:     models.StatusNew,
			Priority:   models.PriorityHigh,
			Deadline:   now,
			AssigneeID: 2,
		},
	}
	facade := NewTaskFacade(
		tasks,
		fakeProjects{exists: map[int64]bool{1: true}},
		fakeUsers{exists: map[int64]bool{2: true}},
		&fakeLogger{},
	)
	facade.SetNow(func() time.Time { return now })

	if err := facade.ChangeStatus(context.Background(), 10, models.StatusDone); err != nil {
		t.Fatalf("expected status change, got error: %v", err)
	}
	if !tasks.statusUpdated {
		t.Fatal("expected status update")
	}
	if tasks.task.Status != models.StatusDone {
		t.Fatalf("expected status %q, got %q", models.StatusDone, tasks.task.Status)
	}
	if len(tasks.task.History) != 1 {
		t.Fatalf("expected one history record, got %d", len(tasks.task.History))
	}
}

func TestTaskFacadeChangeStatusForAssigneeRejectsAnotherUser(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	tasks := &fakeTasks{
		task: models.Task{
			BaseEntity: models.BaseEntity{ID: 10},
			ProjectID:  1,
			Title:      "Write report",
			Status:     models.StatusNew,
			Priority:   models.PriorityHigh,
			Deadline:   now,
			AssigneeID: 2,
		},
	}
	facade := NewTaskFacade(
		tasks,
		fakeProjects{exists: map[int64]bool{1: true}},
		fakeUsers{exists: map[int64]bool{2: true}},
		&fakeLogger{},
	)
	facade.SetNow(func() time.Time { return now })

	if err := facade.ChangeStatusForAssignee(context.Background(), 10, 3, models.StatusDone); err == nil {
		t.Fatal("expected error when task belongs to another user")
	}
	if tasks.statusUpdated {
		t.Fatal("status must not be updated for another user's task")
	}
}
