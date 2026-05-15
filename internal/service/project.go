package service

import (
	"context"
	"fmt"

	"taskflow/internal/models"
)

type ProjectStore interface {
	Create(ctx context.Context, project *models.Project) error
	List(ctx context.Context) ([]models.Project, error)
	ListByOwner(ctx context.Context, ownerID int64) ([]models.Project, error)
	Exists(ctx context.Context, id int64) (bool, error)
	OwnedBy(ctx context.Context, projectID int64, ownerID int64) (bool, error)
	Delete(ctx context.Context, id int64) error
}

type ProjectService struct {
	projects ProjectStore
	users    UserStore
	logger   AppLogger
}

func NewProjectService(projects ProjectStore, users UserStore, logger AppLogger) *ProjectService {
	return &ProjectService{projects: projects, users: users, logger: logger}
}

func (s *ProjectService) Create(ctx context.Context, project models.Project) (models.Project, error) {
	if err := project.Validate(); err != nil {
		s.logger.Error("project validation failed", "error", err)
		return models.Project{}, err
	}
	ownerExists, err := s.users.Exists(ctx, project.OwnerID)
	if err != nil {
		s.logger.Error("failed to check project owner", "error", err)
		return models.Project{}, err
	}
	if !ownerExists {
		err := fmt.Errorf("owner_id must reference an existing user")
		s.logger.Error("project owner does not exist", "owner_id", project.OwnerID)
		return models.Project{}, err
	}
	if err := s.projects.Create(ctx, &project); err != nil {
		s.logger.Error("failed to create project", "error", err)
		return models.Project{}, err
	}
	s.logger.Info("project created", "project_id", project.ID)
	return project, nil
}

func (s *ProjectService) List(ctx context.Context) ([]models.Project, error) {
	return s.projects.List(ctx)
}

func (s *ProjectService) ListByOwner(ctx context.Context, ownerID int64) ([]models.Project, error) {
	return s.projects.ListByOwner(ctx, ownerID)
}

func (s *ProjectService) Delete(ctx context.Context, projectID int64, userID int64) error {
	owned, err := s.projects.OwnedBy(ctx, projectID, userID)
	if err != nil {
		s.logger.Error("failed to check project ownership", "error", err)
		return err
	}
	if !owned {
		err := fmt.Errorf("you can only delete your own projects")
		s.logger.Error("cannot delete project", "project_id", projectID, "user_id", userID)
		return err
	}
	if err := s.projects.Delete(ctx, projectID); err != nil {
		s.logger.Error("failed to delete project", "error", err)
		return err
	}
	s.logger.Info("project deleted", "project_id", projectID)
	return nil
}
