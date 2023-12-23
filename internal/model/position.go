package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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
		CompanyID int    `json:"company_id" db:"company_id"`
		Name      string `json:"name"       db:"name"`
	}

	PositionCourse struct {
		CourseID   int `json:"course_id"`
		PositionID int `json:"position_id"`
	}
)

func (p *PositionSet) Validation() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.Name, validation.Required, validation.Length(2, 256), is.UTFLetterNumeric, validation.NotIn([]rune{'*', '#'})),
		validation.Field(&p.CompanyID, validation.Required),
	)
}
