package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"taskflow/internal/models"
)

type AppLogger interface {
	Info(message string, args ...any)
	Error(message string, args ...any)
}

type UserStore interface {
	Create(ctx context.Context, user *models.User) error
	List(ctx context.Context) ([]models.User, error)
	Exists(ctx context.Context, id int64) (bool, error)
	LoginExists(ctx context.Context, login string) (bool, error)
	FindByID(ctx context.Context, id int64) (models.User, error)
	FindByLogin(ctx context.Context, login string) (models.User, error)
}

type UserService struct {
	users  UserStore
	logger AppLogger
}

func NewUserService(users UserStore, logger AppLogger) *UserService {
	return &UserService{users: users, logger: logger}
}

func (s *UserService) Create(ctx context.Context, user models.User) (models.User, error) {
	if err := user.Validate(); err != nil {
		s.logger.Error("user validation failed", "error", err)
		return models.User{}, err
	}
	exists, err := s.users.LoginExists(ctx, user.Login)
	if err != nil {
		s.logger.Error("failed to check login", "error", err)
		return models.User{}, err
	}
	if exists {
		err := fmt.Errorf("login must be unique")
		s.logger.Error("user login already exists", "login", user.Login)
		return models.User{}, err
	}
	if err := s.users.Create(ctx, &user); err != nil {
		s.logger.Error("failed to create user", "error", err)
		return models.User{}, err
	}
	s.logger.Info("user created", "user_id", user.ID)
	return user, nil
}

func (s *UserService) Register(ctx context.Context, login string, password string) (models.User, error) {
	if err := validatePassword(password); err != nil {
		s.logger.Error("password validation failed", "error", err)
		return models.User{}, err
	}
	hash, err := hashPassword(password)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return models.User{}, err
	}
	return s.Create(ctx, models.User{
		Login:        strings.TrimSpace(login),
		PasswordHash: hash,
	})
}

func (s *UserService) Login(ctx context.Context, login string, password string) (models.User, error) {
	user, err := s.users.FindByLogin(ctx, strings.TrimSpace(login))
	if err != nil {
		s.logger.Error("failed to find user by login", "error", err)
		return models.User{}, fmt.Errorf("invalid login or password")
	}
	if !verifyPassword(password, user.PasswordHash) {
		s.logger.Error("invalid user password", "login", login)
		return models.User{}, fmt.Errorf("invalid login or password")
	}
	s.logger.Info("user logged in", "user_id", user.ID)
	return user, nil
}

func (s *UserService) FindByID(ctx context.Context, id int64) (models.User, error) {
	return s.users.FindByID(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	return s.users.List(ctx)
}

func validatePassword(password string) error {
	if len([]rune(password)) < 6 || len([]rune(password)) > 72 {
		return fmt.Errorf("password length must be from 6 to 72 characters")
	}
	return nil
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	sum := sha256.Sum256(append(salt, []byte(password)...))
	return base64.RawStdEncoding.EncodeToString(salt) + ":" + base64.RawStdEncoding.EncodeToString(sum[:]), nil
}

func verifyPassword(password string, stored string) bool {
	parts := strings.Split(stored, ":")
	if len(parts) != 2 {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	actual := sha256.Sum256(append(salt, []byte(password)...))
	return hmac.Equal(expected, actual[:])
}
