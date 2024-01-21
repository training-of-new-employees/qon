package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/training-of-new-employees/qon/internal/errs"
)

const (
	minNameL = 5
	maxNameL = 256
	minDescL = 10
	maxDescL = 512
)

type Course struct {
	ID          int       `db:"id" json:"id"`
	CreatedBy   int       `db:"created_by" json:"created_by"`
	IsActive    bool      `db:"active" json:"active"`
	IsArchived  bool      `db:"archived" json:"archived"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Status      string    `json:"status,omitempty"`
}

type CoursePreview struct {
	CourseID    int    `json:"course_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type CourseSet struct {
	ID          int    `db:"id" json:"-"`
	CreatedBy   int    `db:"created_by" json:"-"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description,omitempty"`
	IsArchived  bool   `db:"archived" json:"archived,omitempty"`
}

func (c *Course) Validation() error {
	err := validation.Validate(&c.Name, validation.Required)
	if err != nil {
		return errs.ErrCourseNameIsEmpty
	}
	err = validation.Validate(c.Name,
		validation.RuneLength(minNameL, maxNameL),
		validation.By(validateCourseName(c.Name)))
	if err != nil {
		return errs.ErrCourseNameInvalid
	}
	err = validation.Validate(c.Description,
		validation.RuneLength(minDescL, maxDescL),
		validation.By(validateCourseName(c.Description)))
	if err != nil && err != errSpaceEmpty {
		return errs.ErrCourseDescriptionInvalid
	}
	return nil
}

func (cs *CourseSet) Validation() error {
	c := Course{
		Name:        cs.Name,
		Description: cs.Description,
	}
	return c.Validation()
}
func NewCourseSet(id int, creator int) CourseSet {
	return CourseSet{
		ID:        id,
		CreatedBy: creator,
	}
}
