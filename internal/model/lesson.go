package model

type (
	Lesson struct {
		ID         int    `db:"id"          json:"id"`
		CourseID   int    `db:"course_id"   json:"course_id"`
		Name       string `db:"name"        json:"name"`
		Content    string `db:"content"     json:"content"`
		URLPicture string `db:"url_picture" json:"url_picture"`
		Archived   bool   `db:"archived"    json:"archived"`
	}
	/*	LessonCreate struct {
		ID         int    `db:"id"          json:"id"`
		CourseID   int    `db:"course_id"   json:"course_id"`
		Name       string `db:"name"        json:"name"`
		Content    string `db:"content"     json:"content"`
		URLPicture string `db:"url_picture" json:"url_picture"`
	}*/
	LessonUpdate struct {
		ID         int    `db:"id"          json:"-"`
		Name       string `db:"name"        json:"name"`
		Content    string `db:"content"     json:"content"`
		URLPicture string `db:"url_picture" json:"url_picture"`
		Archived   bool   `db:"archived"    json:"archived"`
	}
)
