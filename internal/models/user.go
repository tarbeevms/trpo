package models

import (
	"fmt"
	"strings"
)

type User struct {
	BaseEntity
	AuditInfo
	SoftDelete
	Login        string
	PasswordHash string
}

func (u User) Validate() error {
	if err := validateLength("login", u.Login, 3, 50); err != nil {
		return err
	}
	if strings.Contains(strings.TrimSpace(u.Login), " ") {
		return fmt.Errorf("login must not contain spaces")
	}
	return nil
}
