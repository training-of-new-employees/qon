package rest

type reqCreateUser struct {
	PositionID int    `json:"position_id" db:"position_id"`
	Email      string `json:"email" db:"email"`
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic" db:"patronymic"`
}

type reqCreateLesson struct {
	CourseID   int    `db:"course_id"   json:"course_id"`
	Name       string `db:"name"        json:"name"`
	Content    string `db:"content"     json:"content"`
	URLPicture string `db:"url_picture" json:"url_picture"`
}
