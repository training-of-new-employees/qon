package model

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/training-of-new-employees/qon/internal/errs"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID         int       `db:"id"           json:"id"`
		CompanyID  int       `db:"company_id"   json:"company_id"`
		PositionID int       `db:"position_id"  json:"position_id"`
		Email      string    `db:"email"        json:"email"`
		Password   string    `db:"enc_password" json:"-"`
		IsActive   bool      `db:"active"       json:"active"`
		IsAdmin    bool      `db:"admin"        json:"admin"`
		Name       string    `db:"name"         json:"name"`
		Surname    string    `db:"surname"      json:"surname"`
		Patronymic string    `db:"patronymic"   json:"patronymic"`
		CreatedAt  time.Time `db:"created_at"   json:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"   json:"updated_at"`
	}

	UserSignIn struct {
		Email    string `json:"email"    db:"email"`
		Password string `json:"password" db:"password"`
	}

	UserCreate struct {
		CompanyID  int    `json:"company_id" db:"company_id"`
		PositionID int    `json:"position_id" db:"position_id"`
		Email      string `json:"email" db:"email"`
		Password   string `json:"password" db:"enc_password"`
		IsActive   bool   `json:"active" db:"active"`
		IsAdmin    bool   `json:"admin" db:"admin"`
		Name       string `json:"name" db:"name"`
		Surname    string `json:"surname" db:"surname"`
		Patronymic string `json:"patronymic" db:"patronymic"`
	}

	UserEdit struct {
		ID         int     `json:"-" db:"id"`
		IsArchived bool    `json:"-" db:"archived"`
		Name       *string `json:"name,omitempty" db:"name"`
		Surname    *string `json:"surname,omitempty" db:"surname"`
		Patronymic *string `json:"patronymic,omitempty" db:"patronymic"`
		PositionID *int    `json:"position_id,omitempty" db:"position_id"`
		Email      *string `json:"email,omitempty" db:"email"`
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

func (u *UserCreate) Validation() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.CompanyID, validation.Required),
		validation.Field(&u.PositionID, validation.Required),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Surname, validation.Required),
	)
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

	if u.Password == "" || u.Email == "" {
		return fmt.Errorf("error password or email is empty")
	}

	return nil
}

type (
	CreateAdmin struct {
		Email    string `json:"email"        db:"email"`
		Password string `json:"password"     db:"enc_password"`
		Company  string `json:"company_name" db:"name"`
	}

	AdminCreate struct {
		CompanyID  int    `json:"company_id"  db:"company_id"`
		PositionID int    `json:"position_id" db:"position_id"`
		Email      string `json:"email"       db:"email"`
		Password   string `json:"password"    db:"enc_password"`
		IsActive   bool   `json:"active"      db:"active"`
		IsAdmin    bool   `json:"admin"       db:"admin"`
		Name       string `json:"name"        db:"name"`
		Surname    string `json:"surname"     db:"surname"`
		Patronymic string `json:"patronymic"  db:"patronymic"`
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

func NewAdminCreate(email string, password string) AdminCreate {

	return AdminCreate{
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

func (u *CreateAdmin) Validation() error {
	if err := validation.Validate(&u.Email, validation.Required); err != nil {
		return errs.ErrEmailNotEmpty
	}

	if err := validation.Validate(&u.Email, is.Email); err != nil {
		return errs.ErrInvalidEmail
	}

	if err := validation.Validate(&u.Company, validation.Required); err != nil {
		return errs.ErrCompanyNameNotEmpty
	}

	if err := validation.Validate(&u.Company, validation.Length(1, 256)); err != nil {
		return errs.ErrIncorrectCompanyName
	}

	if err := validation.Validate(&u.Company, is.UTFLetterNumeric); err != nil {
		return errs.ErrIncorrectCompanyName
	}

	if err := validation.Validate(&u.Company, validation.NotIn([]rune{'*', '#'})); err != nil {
		return errs.ErrIncorrectCompanyName
	}

	if err := validation.Validate(&u.Password, validation.Required); err != nil {
		return errs.ErrPasswordNotEmpty
	}

	if err := validation.Validate(&u.Password, validation.Length(6, 30)); err != nil {
		return errs.ErrInvalidPassword
	}

	if err := validation.Validate(&u.Password, validation.By(validatePassword(u.Password))); err != nil {
		return err
	}

	return nil
}

func validatePassword(password string) validation.RuleFunc {
	return func(value interface{}) error {
		if ok := regexp.MustCompile(`\d`).MatchString(password); !ok {
			return errs.ErrInvalidPassword
		}
		if ok := regexp.MustCompile(`[a-z]`).MatchString(password); !ok {
			return errs.ErrInvalidPassword
		}
		if ok := regexp.MustCompile(`[A-Z]`).MatchString(password); !ok {
			return errs.ErrInvalidPassword
		}
		if ok := strings.ContainsAny(password, "!@#$%^&*()_+"); !ok {
			return errs.ErrInvalidPassword
		}
		return nil
	}
}

func (u *CreateAdmin) SetPassword() error {

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return fmt.Errorf("error SetPassword: %v", err)
	}

	u.Password = string(hash)

	return nil

}

func (e *AdminEdit) Validation() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Email, is.Email, validation.Length(5, 50)),
		validation.Field(&e.Company, validation.Length(3, 30)),
		validation.Field(&e.Name, validation.Length(0, 128)),
		validation.Field(&e.Surname, validation.Length(0, 128)),
		validation.Field(&e.Patronymic, validation.Length(0, 128)),
	)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GeneratePassword() string {
	password := make([]byte, 6)
	for i := 0; i < 6; i++ {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}

func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
