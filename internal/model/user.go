package model

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/training-of-new-employees/qon/internal/errs"
)

type (
	User struct {
		ID           int       `db:"id"           json:"id"`
		CompanyID    int       `db:"company_id"   json:"company_id"`
		PositionID   int       `db:"position_id"  json:"position_id"`
		PositionName string    `db:"position_name"  json:"position_name"`
		Email        string    `db:"email"        json:"email"`
		Password     string    `db:"enc_password" json:"-"`
		IsActive     bool      `db:"active"       json:"active"`
		IsArchived   bool      `db:"archived"     json:"archived"`
		IsAdmin      bool      `db:"admin"        json:"admin"`
		Name         string    `db:"name"         json:"name"`
		Surname      string    `db:"surname"      json:"surname"`
		Patronymic   string    `db:"patronymic"   json:"patronymic"`
		CreatedAt    time.Time `db:"created_at"   json:"created_at"`
		UpdatedAt    time.Time `db:"updated_at"   json:"updated_at"`
	}

	UserSignIn struct {
		Email    string `json:"email"    db:"email"`
		Password string `json:"password" db:"password"`
	}

	UserActivation struct {
		Email    string `json:"email"    db:"email"`
		Password string `json:"password" db:"password"`
		Invite   string `json:"invite"`
	}

	UserCreate struct {
		CompanyID  int    `json:"company_id" db:"company_id"`
		PositionID int    `json:"position_id" db:"position_id"`
		IsActive   bool   `json:"active" db:"active"`
		IsArchived bool   `json:"archived" db:"archived"`
		IsAdmin    bool   `json:"admin" db:"admin"`
		Email      string `json:"email" db:"email"`
		Password   string `json:"password" db:"enc_password"`
		Name       string `json:"name" db:"name"`
		Surname    string `json:"surname" db:"surname"`
		Patronymic string `json:"patronymic" db:"patronymic"`
	}

	UserEdit struct {
		ID         int     `json:"-" db:"id"`
		CompanyID  *int    `json:"company_id,omitempty" db:"company_id"`
		PositionID *int    `json:"position_id,omitempty" db:"position_id"`
		IsActive   *bool   `json:"active" db:"active"`
		IsArchived *bool   `json:"archived" db:"archived"`
		Email      *string `json:"email,omitempty" db:"email"`
		Name       *string `json:"name,omitempty" db:"name"`
		Patronymic *string `json:"patronymic,omitempty" db:"patronymic"`
		Surname    *string `json:"surname,omitempty" db:"surname"`
	}

	UserInfo struct {
		User
		CompanyName  string `json:"company_name" db:"company_name"`
		PositionName string `json:"position_name" db:"position_name"`
	}
	EmailReset struct {
		Email string `json:"email"`
	}
)

// Validation - валидация входящих данных при регистрации сотрудника.
func (u *UserCreate) Validation() error {
	// проверка на пустоту поля емейл
	if err := validation.Validate(&u.Email, validation.Required); err != nil {
		return errs.ErrEmailNotEmpty
	}
	// Проверка емейла на корректность
	if err := validation.Validate(&u.Email, validation.By(validateEmail(&u.Email))); err != nil {
		return errs.ErrInvalidEmail
	}
	// проверка на пустоту id компании
	if err := validation.Validate(&u.CompanyID, validation.Required); err != nil {
		return errs.ErrCompanyIDNotEmpty
	}
	// проверка на пустоту id должности
	if err := validation.Validate(&u.PositionID, validation.Required); err != nil {
		return errs.ErrPositionIDNotEmpty
	}
	// проверка на пустоту имени пользователя
	if err := validation.Validate(&u.Name, validation.Required); err != nil {
		return errs.ErrUserNameNotEmpty
	}
	// проверка требований к имени
	if err := validation.Validate(&u.Name, validation.RuneLength(2, 128), validation.By(validateUserName(&u.Name))); err != nil {
		return errs.ErrInvalidUserName
	}
	// проверка на пустоту фамилии пользователя
	if err := validation.Validate(&u.Surname, validation.Required); err != nil {
		return errs.ErrUserSurnameNotEmpty
	}
	// проверка требований к фамилии
	if err := validation.Validate(&u.Surname, validation.RuneLength(2, 128), validation.By(validateUserName(&u.Surname))); err != nil {
		return errs.ErrInvalidUserSurname
	}
	if u.Patronymic != "" {
		// проверка требований к отчеству
		if err := validation.Validate(&u.Patronymic, validation.RuneLength(2, 128), validation.By(validateUserName(&u.Patronymic))); err != nil {
			return errs.ErrInvalidUserPatronymic
		}
	}

	return nil
}

func (ue *UserEdit) Validation() error {
	if ue.Email != nil {
		// проверка на пустоту поля емейл
		if err := validation.Validate(&ue.Email, validation.Required); err != nil {
			return errs.ErrEmailNotEmpty
		}
		// Проверка емейла на корректность
		if err := validation.Validate(&ue.Email, validation.By(validateEmail(ue.Email))); err != nil {
			return errs.ErrInvalidEmail
		}
	}
	// проверка требований к имени
	if ue.Name != nil {
		if err := validation.Validate(&ue.Name, validation.Required); err != nil {
			return errs.ErrUserNameNotEmpty
		}
		if err := validation.Validate(&ue.Name, validation.RuneLength(2, 128), validation.By(validateUserName(ue.Name))); err != nil {
			return errs.ErrInvalidUserName
		}
	}
	if ue.Surname != nil {
		if err := validation.Validate(&ue.Surname, validation.Required); err != nil {
			return errs.ErrUserSurnameNotEmpty
		}
		// проверка требований к фамилии
		if err := validation.Validate(&ue.Surname, validation.RuneLength(2, 128), validation.By(validateUserName(ue.Surname))); err != nil {
			return errs.ErrInvalidUserSurname
		}
	}

	if ue.Patronymic != nil {
		// проверка требований к отчеству
		if err := validation.Validate(&ue.Patronymic, validation.RuneLength(2, 128), validation.By(validateUserName(ue.Patronymic))); err != nil {
			return errs.ErrInvalidUserPatronymic
		}
	}
	return nil
}

func (u *UserCreate) SetPassword() error {

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return fmt.Errorf("error SetPassword: %v", err)
	}

	u.Password = string(hash)

	return nil
}

func (u *UserCreate) SetActive() error {

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return fmt.Errorf("error SetPassword: %v", err)
	}

	u.Password = string(hash)

	return nil
}

func (u *User) CheckPassword(password string) error {

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return fmt.Errorf("error hash password validation: %v", err)
	}

	return nil
}

func (u *UserSignIn) Validation() error {
	if u.Email == "" || u.Password == "" {
		return errs.ErrEmailOrPasswordEmpty
	}

	var err error
	u.Email, err = modifyEmail(u.Email)
	if err != nil {
		return errs.ErrIncorrectEmailOrPassword
	}

	return nil
}

func (u *UserActivation) Validation() error {
	if u.Email == "" || u.Password == "" {
		return errs.ErrEmailOrPasswordEmpty
	}

	if err := validation.Validate(&u.Email, validation.Required); err != nil {
		return errs.ErrEmailNotEmpty
	}

	// Проверка емейла на корректность
	if err := validation.Validate(&u.Email, validation.By(validateEmail(&u.Email))); err != nil {
		return errs.ErrInvalidEmail
	}

	if err := validation.Validate(&u.Password, validation.Required); err != nil {
		return errs.ErrPasswordNotEmpty
	}

	if err := validation.Validate(&u.Password, validation.Length(6, 30)); err != nil {
		return errs.ErrInvalidPassword
	}

	if err := validation.Validate(&u.Password, validation.By(validatePassword(u.Password))); err != nil {
		return errs.ErrInvalidPassword
	}

	if err := validation.Validate(&u.Invite, validation.Required); err != nil {
		return errs.ErrInviteNotEmpty
	}

	return nil
}

type (
	CreateAdmin struct {
		Email    string `json:"email"        db:"email"`
		Password string `json:"password"     db:"enc_password"`
		Company  string `json:"company_name" db:"name"`
	}

	// AdminEdit - Структура для передачи изменяемых данных администратора
	AdminEdit struct {
		ID         int     `json:"id,omitempty"           db:"id"`
		Email      *string `json:"email,omitempty"        db:"email"`
		Company    *string `json:"company_name,omitempty" db:"company_name"`
		Name       *string `json:"name,omitempty"         db:"name"`
		Patronymic *string `json:"patronymic,omitempty"   db:"patronymic"`
		Surname    *string `json:"surname,omitempty"      db:"surname"`
	}
)

func NewAdminCreate(email string, password string) UserCreate {
	return UserCreate{
		CompanyID:  0,
		PositionID: 0,
		Email:      email,
		Password:   password,
		IsActive:   true,
		IsAdmin:    true,
		Name:       "admin",
		Surname:    "admin",
		Patronymic: "admin",
	}
}

// Validation - валидация входящих данных при регистрации администратора.
func (u *CreateAdmin) Validation() error {
	// Проверка на пустоту поля емейл
	if err := validation.Validate(&u.Email, validation.Required); err != nil {
		return errs.ErrEmailNotEmpty
	}
	// Проверка емейла на корректность
	if err := validation.Validate(&u.Email, validation.By(validateEmail(&u.Email))); err != nil {
		return errs.ErrInvalidEmail
	}
	// Проверка на пустоту пароля
	if err := validation.Validate(&u.Password, validation.Required); err != nil {
		return errs.ErrPasswordNotEmpty
	}
	// Проверка на длину пароля
	if err := validation.Validate(&u.Password, validation.Length(6, 30)); err != nil {
		return errs.ErrInvalidPassword
	}
	// Проверка состава пароля
	if err := validation.Validate(&u.Password, validation.By(validatePassword(u.Password))); err != nil {
		return errs.ErrInvalidPassword
	}
	// Проверка на пустоту имени компании
	if err := validation.Validate(&u.Company, validation.Required); err != nil {
		return errs.ErrCompanyNameNotEmpty
	}
	// Проверка имени компании на состав
	if err := validation.Validate(&u.Company, validation.Length(1, 256), validation.By(validateNameDescription(&u.Company))); err != nil {
		return errs.ErrInvalidCompanyName
	}

	return nil
}

func (u *CreateAdmin) SetPassword() error {

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return fmt.Errorf("error SetPassword: %v", err)
	}

	u.Password = string(hash)

	return nil

}

func (ae *AdminEdit) Validation() error {
	if ae.Email != nil {
		// проверка на пустоту поля емейл
		if err := validation.Validate(&ae.Email, validation.Required); err != nil {
			return errs.ErrEmailNotEmpty
		}
		// Проверка емейла на корректность
		if err := validation.Validate(&ae.Email, validation.By(validateEmail(ae.Email))); err != nil {
			return errs.ErrInvalidEmail
		}
	}
	// проверка требований к имени
	if ae.Name != nil {
		if err := validation.Validate(&ae.Name, validation.Required); err != nil {
			return errs.ErrUserNameNotEmpty
		}
		if err := validation.Validate(&ae.Name, validation.RuneLength(2, 128), validation.By(validateUserName(ae.Name))); err != nil {
			return errs.ErrInvalidUserName
		}
	}
	if ae.Surname != nil {
		if err := validation.Validate(&ae.Surname, validation.Required); err != nil {
			return errs.ErrUserSurnameNotEmpty
		}
		// проверка требований к фамилии
		if err := validation.Validate(&ae.Surname, validation.RuneLength(2, 128), validation.By(validateUserName(ae.Surname))); err != nil {
			return errs.ErrInvalidUserSurname
		}
	}

	if ae.Patronymic != nil {
		// проверка требований к отчеству
		if err := validation.Validate(&ae.Patronymic, validation.RuneLength(2, 128), validation.By(validateUserName(ae.Patronymic))); err != nil {
			return errs.ErrInvalidUserPatronymic
		}
	}

	if ae.Company != nil {
		// Проверка имени компании на состав
		if err := validation.Validate(&ae.Company, validation.Length(1, 256), validation.By(validateNameDescription(ae.Company))); err != nil {
			return errs.ErrInvalidCompanyName
		}
	}

	return nil
}

func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
