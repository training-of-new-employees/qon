package rest

type reqCreatePosition struct {
	CompanyID int    `json:"company_id" example:"1" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

type reqCreateUser struct {
	PositionID int    `json:"position_id" example:"1" binding:"required"`
	Email      string `json:"email" example:"somebody@example.org" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
}

type reqEditUser struct {
	Email      *string `json:"email" example:"somebody@example.org"`
	Name       *string `json:"name"`
	Patronymic *string `json:"patronymic"`
	Surname    *string `json:"surname"`
	IsArchived *bool   `json:"archived"`
	IsActive   *bool   `json:"active"`
}

type reqEditAdmin struct {
	Email      *string `json:"email" example:"somebody@example.org"`
	Company    *string `json:"company_name"`
	Name       *string `json:"name"`
	Patronymic *string `json:"patronymic"`
	Surname    *string `json:"surname"`
}

type reqCreateCourse struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type reqCreateLesson struct {
	CourseID   int    `json:"course_id" example:"1" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Content    string `json:"content" binding:"required"`
	URLPicture string `json:"url_picture"`
}
