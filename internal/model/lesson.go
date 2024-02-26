package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/training-of-new-employees/qon/internal/errs"
)

const (
	minLessonNameL = 1
	maxLessonNameL = 256
	minContentL    = 20
	maxContentL    = 65000
	minURLPictureL = 5
	maxURLPictureL = 1024
)

type (
	Lesson struct {
		ID         int    `db:"id"          json:"id"`
		CourseID   int    `db:"course_id"   json:"course_id"`
		Name       string `db:"name"        json:"name"`
		Content    string `db:"content"     json:"content"`
		URLPicture string `db:"url_picture" json:"url_picture"`
		Archived   bool   `db:"archived"    json:"archived"`
		Status     string `json:"status,omitempty"`
	}

	LessonPreview struct {
		LessonID int    `json:"lesson_id"`
		CourseID int    `json:"course_id"`
		Name     string `json:"name"`
		Status   string `json:"status"`
	}
	LessonUpdate struct {
		ID         int    `db:"id"          json:"-"`
		Name       string `db:"name"        json:"name"`
		Content    string `db:"content"     json:"content"`
		URLPicture string `db:"url_picture" json:"url_picture"`
		Archived   bool   `db:"archived"    json:"archived"`
	}
	LessonStatusUpdate struct {
		Status string `json:"status"`
	}
)

// Validation - валидация данных при создании урока.
func (l *Lesson) Validation() error {
	if err := validation.Validate(&l.Name, validation.Required); err != nil {
		return errs.ErrLessonNameNotEmpty
	}

	if err := validation.Validate(&l.Name, validation.RuneLength(minLessonNameL, maxLessonNameL), validation.By(validateObjName(&l.Name))); err != nil {
		return errs.ErrInvalidLessonName
	}

	if err := validation.Validate(&l.Content, validation.Required); err != nil {
		return errs.ErrTextContentNotEmpty
	}

	if err := validation.Validate(&l.Content, validation.RuneLength(minContentL, maxContentL), validation.By(validateObjDescription(&l.Content))); err != nil {
		return errs.ErrInvalidTextContent
	}

	if l.URLPicture != "" {
		if err := validation.Validate(&l.URLPicture, validation.RuneLength(minURLPictureL, maxURLPictureL), validation.By(validateURLPicture(&l.URLPicture))); err != nil {
			return errs.ErrInvalidURLPicture
		}
	}

	return nil
}

// Validation - валидация данных при редактировании урока.
func (l *LessonUpdate) Validation() error {
	if l.Name != "" {
		if err := validation.Validate(&l.Name, validation.RuneLength(minLessonNameL, maxLessonNameL), validation.By(validateObjName(&l.Name))); err != nil {
			return errs.ErrInvalidLessonName
		}
	}

	if l.Content != "" {
		if err := validation.Validate(&l.Content, validation.RuneLength(minContentL, maxContentL), validation.By(validateObjDescription(&l.Content))); err != nil {
			return errs.ErrInvalidTextContent
		}
	}

	if l.URLPicture != "" {
		if err := validation.Validate(&l.URLPicture, validation.RuneLength(minURLPictureL, maxURLPictureL), validation.By(validateURLPicture(&l.URLPicture))); err != nil {
			return errs.ErrInvalidURLPicture
		}
	}

	return nil
}
