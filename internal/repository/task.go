package repository

import (
	"context"
	"database/sql"
	"strconv"

	"taskflow/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
		INSERT INTO tasks (project_id, title, description, status, priority, deadline, assignee_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`, task.ProjectID, task.Title, task.Description, task.Status, task.Priority, task.Deadline, task.AssigneeID, task.CreatedBy).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return err
	}

	for _, tag := range task.Tags {
		var tagID int64
		if err := tx.QueryRowContext(ctx, `
			INSERT INTO tags (name)
			VALUES ($1)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, tag.Name).Scan(&tagID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO task_tags (task_id, tag_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, task.ID, tagID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *TaskRepository) List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
	query := `
		SELECT id, project_id, title, description, status, priority, deadline, assignee_id, created_at, updated_at, created_by
		FROM tasks
		WHERE deleted_at IS NULL
	`
	var args []any
	if filter.AssigneeID > 0 {
		args = append(args, filter.AssigneeID)
		query += ` AND assignee_id = $` + intPlaceholder(len(args))
	}
	query += ` ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.Deadline,
			&task.AssigneeID,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.CreatedBy,
		); err != nil {
			return nil, err
		}
		tags, err := r.listTags(ctx, task.ID)
		if err != nil {
			return nil, err
		}
		task.Tags = tags
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) FindByID(ctx context.Context, id int64) (models.Task, error) {
	var task models.Task
	err := r.db.QueryRowContext(ctx, `
		SELECT id, project_id, title, description, status, priority, deadline, assignee_id, created_at, updated_at, created_by
		FROM tasks
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.Deadline,
		&task.AssigneeID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.CreatedBy,
	)
	if err != nil {
		return models.Task{}, err
	}
	task.Tags, err = r.listTags(ctx, task.ID)
	return task, err
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, task models.Task, history models.TaskHistory) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		UPDATE tasks
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`, task.Status, task.ID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO task_history (task_id, old_status, new_status, changed_at)
		VALUES ($1, $2, $3, $4)
	`, history.TaskID, history.OldStatus, history.NewStatus, history.ChangedAt); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *TaskRepository) listTags(ctx context.Context, taskID int64) ([]models.Tag, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT tags.id, tags.name
		FROM tags
		JOIN task_tags ON task_tags.tag_id = tags.id
		WHERE task_tags.task_id = $1
		ORDER BY tags.name
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

// Delete выполняет мягкое удаление задачи (устанавливает deleted_at)
func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tasks
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	return err
}

func intPlaceholder(value int) string {
	return strconv.Itoa(value)
}
