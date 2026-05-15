package service

import (
	"context"
	"testing"

	"taskflow/internal/models"
)

type fakeUserStore struct {
	users map[string]models.User
	next  int64
}

func newFakeUserStore() *fakeUserStore {
	return &fakeUserStore{users: make(map[string]models.User), next: 1}
}

func (f *fakeUserStore) Create(ctx context.Context, user *models.User) error {
	user.ID = f.next
	f.next++
	f.users[user.Login] = *user
	return nil
}

func (f *fakeUserStore) List(ctx context.Context) ([]models.User, error) {
	var users []models.User
	for _, user := range f.users {
		users = append(users, user)
	}
	return users, nil
}

func (f *fakeUserStore) Exists(ctx context.Context, id int64) (bool, error) {
	for _, user := range f.users {
		if user.ID == id {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeUserStore) LoginExists(ctx context.Context, login string) (bool, error) {
	_, ok := f.users[login]
	return ok, nil
}

func (f *fakeUserStore) FindByID(ctx context.Context, id int64) (models.User, error) {
	for _, user := range f.users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, errNotFoundForTest{}
}

func (f *fakeUserStore) FindByLogin(ctx context.Context, login string) (models.User, error) {
	user, ok := f.users[login]
	if !ok {
		return models.User{}, errNotFoundForTest{}
	}
	return user, nil
}

type errNotFoundForTest struct{}

func (errNotFoundForTest) Error() string {
	return "not found"
}

func TestUserServiceRegisterAndLoginPositive(t *testing.T) {
	store := newFakeUserStore()
	service := NewUserService(store, &fakeLogger{})

	registered, err := service.Register(context.Background(), "ivan", "secret1")
	if err != nil {
		t.Fatalf("expected registration, got error: %v", err)
	}
	if registered.ID == 0 {
		t.Fatal("expected registered user id")
	}
	if registered.PasswordHash == "" || registered.PasswordHash == "secret1" {
		t.Fatal("expected password to be hashed")
	}

	loggedIn, err := service.Login(context.Background(), "ivan", "secret1")
	if err != nil {
		t.Fatalf("expected login, got error: %v", err)
	}
	if loggedIn.ID != registered.ID {
		t.Fatalf("expected user id %d, got %d", registered.ID, loggedIn.ID)
	}
}

func TestUserServiceRegisterNegativePassword(t *testing.T) {
	service := NewUserService(newFakeUserStore(), &fakeLogger{})
	if _, err := service.Register(context.Background(), "ivan", "123"); err == nil {
		t.Fatal("expected password validation error")
	}
}

func TestUserServiceLoginNegativePassword(t *testing.T) {
	store := newFakeUserStore()
	service := NewUserService(store, &fakeLogger{})
	if _, err := service.Register(context.Background(), "ivan", "secret1"); err != nil {
		t.Fatalf("expected registration, got error: %v", err)
	}
	if _, err := service.Login(context.Background(), "ivan", "wrong-password"); err == nil {
		t.Fatal("expected login error")
	}
}
