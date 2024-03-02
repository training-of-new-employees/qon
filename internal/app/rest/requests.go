package rest

type reqCreateUser struct {
	PositionID int    `json:"position_id" db:"position_id"`
	Email      string `json:"email" db:"email"`
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic" db:"patronymic"`
	Password   string `json:"password" db:"enc_password"`
}

type reqEditUser struct {
	PositionID *int    `json:"position_id,omitempty" db:"position_id"`
	IsActive   *bool   `json:"active" db:"active"`
	IsArchived *bool   `json:"archived" db:"archived"`
	Email      *string `json:"email,omitempty" db:"email"`
	Name       *string `json:"name,omitempty" db:"name"`
	Patronymic *string `json:"patronymic,omitempty" db:"patronymic"`
	Surname    *string `json:"surname,omitempty" db:"surname"`
}

type reqEditAdmin struct {
	Email      *string `json:"email,omitempty"        db:"email"`
	Company    *string `json:"company_name,omitempty" db:"company_name"`
	Name       *string `json:"name,omitempty"         db:"name"`
	Patronymic *string `json:"patronymic,omitempty"   db:"patronymic"`
	Surname    *string `json:"surname,omitempty"      db:"surname"`
}

type reqCreateLesson struct {
	CourseID   int    `db:"course_id"   json:"course_id"`
	Name       string `db:"name"        json:"name"`
	Content    string `db:"content"     json:"content"`
	URLPicture string `db:"url_picture" json:"url_picture"`
}
