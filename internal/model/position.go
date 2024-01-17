package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/training-of-new-employees/qon/internal/errs"
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

	PositionSet struct {
		CompanyID  int    `json:"company_id" db:"company_id"`
		Name       string `json:"name"       db:"name"`
		IsArchived bool   `json:"archived" db:"archived"`
	}

	PositionCourse struct {
		CourseID   int `json:"course_id"`
		PositionID int `json:"position_id"`
	}

	PositionAssignCourses struct {
		CourseID []int `json:"course_id"`
	}
)

func (p *PositionSet) Validation() error {
	// проверка на пустоту имя компании
	if err := validation.Validate(p.Name, validation.Required); err != nil {
		return errs.ErrPositionNameNotEmpty
	}
	// проверка на корректность имени компании
	if err := validation.Validate(p.Name, validation.RuneLength(2, 256), validation.By(validateCompanyPositionName(p.Name))); err != nil {
		return errs.ErrInvalidPositionName
	}
	// проверка на наличие id компании
	if err := validation.Validate(&p.CompanyID, validation.Required); err != nil {
		return errs.ErrCompanyIDNotEmpty
	}

	return nil
}

func (p *PositionAssignCourses) Validation() error {
	return nil
}
