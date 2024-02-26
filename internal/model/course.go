package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/training-of-new-employees/qon/internal/errs"
)

const (
	minCourseNameL = 1
	maxCourseNameL = 256
	minCourseDescL = 10
	maxCourseDescL = 1024
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

// Validation - валидация данных при создании курса.
func (c *Course) Validation() error {
	if err := validation.Validate(&c.Name, validation.Required); err != nil {
		return errs.ErrCourseNameIsEmpty
	}
	err := validation.Validate(&c.Name, validation.RuneLength(minCourseNameL, maxCourseNameL), validation.By(validateObjName(&c.Name)))
	if err != nil {
		return errs.ErrInvalidCourseName
	}

	if c.Description != "" {
		err = validation.Validate(&c.Description, validation.RuneLength(minCourseDescL, maxCourseDescL), validation.By(validateObjDescription(&c.Description)))
		if err != nil && err != errSpaceEmpty {
			return errs.ErrInvalidCourseDescription
		}
	}

	return nil
}

func NewCourseSet(id int, creator int) CourseSet {
	return CourseSet{
		ID:        id,
		CreatedBy: creator,
	}
}

// Validation - валидация данных курса при редактировании.
func (cs *CourseSet) Validation() error {
	c := Course{
		Name:        cs.Name,
		Description: cs.Description,
	}

	if err := c.Validation(); err != nil {
		return err
	}

	cs.Name = c.Name
	cs.Description = c.Description

	return nil
}
