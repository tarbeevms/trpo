package models

import (
	"fmt"
	"strings"
	"time"
)

type BaseEntity struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetID возвращает ID сущности
func (e BaseEntity) GetID() int64 {
	return e.ID
}

// GetCreatedAt возвращает время создания
func (e BaseEntity) GetCreatedAt() time.Time {
	return e.CreatedAt
}

type AuditInfo struct {
	CreatedBy int64
}

// GetID возвращает ID создателя записи - конфликт с BaseEntity.GetID()
// При множественном наследовании требует явного вызова: task.AuditInfo.GetID()
func (a AuditInfo) GetID() int64 {
	return a.CreatedBy
}

// GetCreatedBy возвращает ID создателя (уникальный метод)
func (a AuditInfo) GetCreatedBy() int64 {
	return a.CreatedBy
}

type SoftDelete struct {
	DeletedAt *time.Time
}

// GetID возвращает 0 если не удалён, или ID записи если удалён
// Ещё один конфликт с BaseEntity.GetID() и AuditInfo.GetID()
func (s SoftDelete) GetID() int64 {
	if s.DeletedAt != nil {
		return -1 // обозначаем удалённую запись
	}
	return 0
}

// IsDeleted проверяет, удалена ли запись (уникальный метод)
func (s SoftDelete) IsDeleted() bool {
	return s.DeletedAt != nil
}

func validateLength(field string, value string, min int, max int) error {
	length := len([]rune(strings.TrimSpace(value)))
	if length < min || length > max {
		return fmt.Errorf("%s length must be from %d to %d characters", field, min, max)
	}
	return nil
}

func validateMaxLength(field string, value string, max int) error {
	length := len([]rune(strings.TrimSpace(value)))
	if length > max {
		return fmt.Errorf("%s length must be up to %d characters", field, max)
	}
	return nil
}

func dateOnly(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}
