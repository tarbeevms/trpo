package repository

import (
	"context"
	"database/sql"

	"taskflow/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO users (login, password_hash, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, user.Login, user.PasswordHash, user.CreatedBy).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, login, password_hash, created_at, updated_at, created_by
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.CreatedBy); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)`, id).Scan(&exists)
	return exists, err
}

func (r *UserRepository) LoginExists(ctx context.Context, login string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1 AND deleted_at IS NULL)`, login).Scan(&exists)
	return exists, err
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, login, password_hash, created_at, updated_at, created_by
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.CreatedBy)
	return user, err
}

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, login, password_hash, created_at, updated_at, created_by
		FROM users
		WHERE login = $1 AND deleted_at IS NULL
	`, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.CreatedBy)
	return user, err
}
