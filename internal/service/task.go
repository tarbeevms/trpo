package service

import (
	"context"
	"fmt"
	"time"

	"taskflow/internal/models"
)

type TaskStore interface {
	Create(ctx context.Context, task *models.Task) error
	List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error)
	FindByID(ctx context.Context, id int64) (models.Task, error)
	UpdateStatus(ctx context.Context, task models.Task, history models.TaskHistory) error
}

type TaskFacade struct {
	tasks    TaskStore
	projects ProjectStore
	users    UserStore
	logger   AppLogger
	now      func() time.Time
}

func NewTaskFacade(tasks TaskStore, projects ProjectStore, users UserStore, logger AppLogger) *TaskFacade {
	return &TaskFacade{
		tasks:    tasks,
		projects: projects,
		users:    users,
		logger:   logger,
		now:      time.Now,
	}
}

func (f *TaskFacade) SetNow(now func() time.Time) {
	f.now = now
}

func (f *TaskFacade) Create(ctx context.Context, task models.Task) (models.Task, error) {
	if task.Status == "" {
		task.Status = models.StatusNew
	}
	if err := task.Validate(f.now()); err != nil {
		f.logger.Error("task validation failed", "error", err)
		return models.Task{}, err
	}
	projectExists, err := f.projects.Exists(ctx, task.ProjectID)
	if err != nil {
		f.logger.Error("failed to check task project", "error", err)
		return models.Task{}, err
	}
	if !projectExists {
		err := fmt.Errorf("project_id must reference an existing project")
		f.logger.Error("task project does not exist", "project_id", task.ProjectID)
		return models.Task{}, err
	}
	if task.CreatedBy > 0 {
		projectOwned, err := f.projects.OwnedBy(ctx, task.ProjectID, task.CreatedBy)
		if err != nil {
			f.logger.Error("failed to check task project owner", "error", err)
			return models.Task{}, err
		}
		if !projectOwned {
			err := fmt.Errorf("project_id must reference current user's project")
			f.logger.Error("task project belongs to another user", "project_id", task.ProjectID)
			return models.Task{}, err
		}
	}
	assigneeExists, err := f.users.Exists(ctx, task.AssigneeID)
	if err != nil {
		f.logger.Error("failed to check task assignee", "error", err)
		return models.Task{}, err
	}
	if !assigneeExists {
		err := fmt.Errorf("assignee_id must reference an existing user")
		f.logger.Error("task assignee does not exist", "assignee_id", task.AssigneeID)
		return models.Task{}, err
	}
	if err := f.tasks.Create(ctx, &task); err != nil {
		f.logger.Error("failed to create task", "error", err)
		return models.Task{}, err
	}
	f.logger.Info("task created", "task_id", task.ID)
	return task, nil
}

func (f *TaskFacade) List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
	return f.tasks.List(ctx, filter)
}

func (f *TaskFacade) ChangeStatus(ctx context.Context, taskID int64, status models.Status) error {
	task, err := f.tasks.FindByID(ctx, taskID)
	if err != nil {
		f.logger.Error("failed to find task", "error", err)
		return err
	}
	history, err := task.ChangeStatus(status, f.now())
	if err != nil {
		f.logger.Error("task status validation failed", "error", err)
		return err
	}
	if err := f.tasks.UpdateStatus(ctx, task, history); err != nil {
		f.logger.Error("failed to update task status", "error", err)
		return err
	}
	f.logger.Info("task status changed", "task_id", task.ID, "status", status)
	return nil
}

func (f *TaskFacade) ChangeStatusForAssignee(ctx context.Context, taskID int64, assigneeID int64, status models.Status) error {
	task, err := f.tasks.FindByID(ctx, taskID)
	if err != nil {
		f.logger.Error("failed to find task", "error", err)
		return err
	}
	if task.AssigneeID != assigneeID {
		err := fmt.Errorf("task must belong to current user")
		f.logger.Error("task belongs to another user", "task_id", taskID)
		return err
	}
	history, err := task.ChangeStatus(status, f.now())
	if err != nil {
		f.logger.Error("task status validation failed", "error", err)
		return err
	}
	if err := f.tasks.UpdateStatus(ctx, task, history); err != nil {
		f.logger.Error("failed to update task status", "error", err)
		return err
	}
	f.logger.Info("task status changed", "task_id", task.ID, "status", status)
	return nil
}
