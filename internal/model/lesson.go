package model

import "time"

type (
	Lesson struct {
		ID          int       `db:"id"          json:"id"`
		CourseID    int       `db:"course_id"   json:"course_id"`
		CreatedBy   int       `db:"created_by"   json:"created_by"`
		Number      int       `db:"number"      json:"number"`
		Name        string    `db:"name"        json:"name"`
		Description string    `db:"description" json:"description"`
		IsActive    bool      `db:"active"      json:"active"`
		IsArchived  bool      `db:"archived"    json:"archived"`
		CreatedAt   time.Time `db:"created_at"  json:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"  json:"updated_at"`
	}
	LessonCreate struct {
		CourseID    int    `db:"course_id"  json:"course_id"`
		Name        string `db:"name"       json:"name"`
		Description string `db:"decription" json:"description"`
		Path        string `db:"path"       json:"path"`
	}
)
