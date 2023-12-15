package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type (
	Course struct {
		ID          int       `db:"id"          json:"id"`
		Name        string    `db:"name"        json:"name"`
		Description string    `db:"description" json:"description"`
		CreatedBy   string    `db:"created_by"  json:"created_by"`
		IsActive    bool      `db:"active"      json:"active"`
		IsArchived  bool      `db:"archived"    json:"archived"`
		CreatedAt   time.Time `db:"created_at" json:"created_at"`
		UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	}

	CourseCreate struct {
		Name        string `json:"name" db:"name"`
		Description string `json:"description" db:"description"`
	}
)

func (c *CourseCreate) Validation() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(3, 255)),
		validation.Field(&c.Description, validation.Required, validation.Length(3, 255)),
	)
}
