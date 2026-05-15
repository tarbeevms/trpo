package models

import "fmt"

type Project struct {
	ID int64
	SoftDelete
	Name        string
	Description string
	OwnerID     int64
	Tasks       []Task
}

func (p Project) Validate() error {
	if err := validateLength("name", p.Name, 3, 80); err != nil {
		return err
	}
	if err := validateMaxLength("description", p.Description, 500); err != nil {
		return err
	}
	if p.OwnerID <= 0 {
		return fmt.Errorf("owner_id is required")
	}
	return nil
}
