package repository

import (
	"context"
	"database/sql"

	"taskflow/internal/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO projects (name, description, owner_id, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, project.Name, project.Description, project.OwnerID, project.CreatedBy).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
}

func (r *ProjectRepository) List(ctx context.Context) ([]models.Project, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, owner_id, created_at, updated_at, created_by
		FROM projects
		WHERE deleted_at IS NULL
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt, &project.CreatedBy); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, rows.Err()
}

func (r *ProjectRepository) ListByOwner(ctx context.Context, ownerID int64) ([]models.Project, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, owner_id, created_at, updated_at, created_by
		FROM projects
		WHERE owner_id = $1 AND deleted_at IS NULL
		ORDER BY id
	`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt, &project.CreatedBy); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, rows.Err()
}

func (r *ProjectRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND deleted_at IS NULL)`, id).Scan(&exists)
	return exists, err
}

func (r *ProjectRepository) OwnedBy(ctx context.Context, projectID int64, ownerID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM projects
			WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL
		)
	`, projectID, ownerID).Scan(&exists)
	return exists, err
}
