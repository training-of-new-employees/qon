package model

import "time"

type (
	Lesson struct {
		ID         int       `db:"id"          json:"id"`
		CourseID   int       `db:"course_id"   json:"course_id"`
		CreatedBy  int       `db:"created_by"  json:"created_by"`
		Number     int       `db:"number"      json:"number"`
		Name       string    `db:"name"        json:"name"`
		Content    string    `db:"content"     json:"content"`
		URLPicture string    `db:"url_picture" json:"url_picture"`
		IsActive   bool      `db:"active"      json:"active"`
		IsArchived bool      `db:"archived"    json:"archived"`
		CreatedAt  time.Time `db:"created_at"  json:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"  json:"updated_at"`
	}
	LessonCreate struct {
		CourseID   int    `db:"course_id"   json:"course_id"`
		Name       string `db:"name"        json:"name"`
		Content    string `db:"content"     json:"content"`
		URLPicture string `db:"url_picture" json:"url_picture"`
	}
	LessonUpdate struct {
		ID         int    `db:"id"         json:"id"`
		CourseID   int    `db:"course_id"  json:"course_id"`
		Name       string `db:"name"       json:"name"`
		Content    string `db:"content"     json:"content"`
		URLPicture string `db:"url_picture" json:"url_picture"`
	}
)
