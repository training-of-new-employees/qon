package model

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

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
		return errs.ErrCourseNameNotEmpty
	}
	err = validation.Validate(&c.Name, validation.Length(minNameL, maxNameL))
	if err != nil {
		return errs.ErrCourseNameInvalid
	}
	nameWOSpaces := strings.ReplaceAll(c.Name, " ", "")
	err = validation.Validate(&nameWOSpaces, is.UTFLetterNumeric, validation.NotIn([]rune{'*', '#'}))
	if err != nil {
		return errs.ErrCourseNameInvalid
	}

	err = validation.Validate(&c.Description, validation.Length(minDescL, maxDescL))
	if err != nil {
		return errs.ErrCourseDescriptionInvalid
	}
	descWOSpaces := strings.ReplaceAll(c.Description, " ", "")
	err = validation.Validate(&descWOSpaces, is.UTFLetterNumeric, validation.NotIn([]rune{'*', '#'}))
	if err != nil {
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
