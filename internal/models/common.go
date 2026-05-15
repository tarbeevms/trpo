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

type AuditInfo struct {
	CreatedBy int64
}

type SoftDelete struct {
	DeletedAt *time.Time
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
