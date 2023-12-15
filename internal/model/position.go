package model

import (
	"fmt"
	"time"
)

type (
	Position struct {
		ID         int       `db:"id"         json:"id"`
		CompanyID  int       `db:"company_id" json:"company_id"`
		Name       string    `db:"name"       json:"name"`
		IsActive   bool      `db:"active"     json:"active"`
		IsArchived bool      `db:"archived"   json:"archived"`
		CreatedAt  time.Time `db:"created_at" json:"created_at"`
		UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	}

	PositionCreate struct {
		CompanyID int    `json:"company_id" db:"company_id"`
		Name      string `json:"name"       db:"name"`
	}

	PositionUpdate struct {
		CompanyID int    `json:"company_id" db:"company_id"`
		Name      string `json:"name"       db:"name"`
	}
	PositionCourse struct {
		CourseID   int `json:"course_id"`
		PositionID int `json:"position_id"`
	}
)

func (p *PositionCreate) Validation() error {
	if p.CompanyID == 0 || p.Name == "" {
		return fmt.Errorf("validation empty values")
	}

	return nil
}

func (p *PositionUpdate) Validation() error {
	if p.CompanyID == 0 || p.Name == "" {
		return fmt.Errorf("validation empty values")
	}

	return nil
}
