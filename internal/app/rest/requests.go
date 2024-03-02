package rest

type reqCreateUser struct {
	PositionID int    `json:"position_id" db:"position_id"`
	Email      string `json:"email" db:"email"`
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic" db:"patronymic"`
}
