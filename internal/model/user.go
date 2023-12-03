package model

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type (
	User struct {
		ID         int       `db:"id" json:"id"`
		LeaderID   int       `db:"leader_id" json:"leader_id"`
		CompanyID  int       `db:"company_id" json:"company_id"`
		PositionID int       `db:"position_id" json:"position_id"`
		Email      string    `db:"email" json:"email"`
		Password   string    `db:"password" json:"password"`
		IsActive   bool      `db:"is_active" json:"is_active"`
		IsAdmin    bool      `db:"is_admin" json:"is_admin"`
		IsDeleted  bool      `db:"is_deleted" json:"is_deleted"`
		CreatedAt  time.Time `db:"created_at" json:"created_at"`
		UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	}

	UserSignIn struct {
		Email    string `json:"email" db:"email"`
		Password string `json:"password" db:"password"`
	}
)

type (
	CreateAdmin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Company  string `json:"company_name"`
	}
)

func (u *CreateAdmin) Validation() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Company, validation.Required, validation.Length(3, 30)))
}

func (u *CreateAdmin) ValidatePassword() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Password, validation.Required, validation.Length(6, 30)))
}

func (u *CreateAdmin) SetPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return fmt.Errorf("err SetPassword: %v", err)
	}

	u.Password = string(hash)

	return nil

}

func (u *CreateAdmin) CheckPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return fmt.Errorf("err hash password validation: %v", err)
	}

	return nil
}

func (u *UserSignIn) Validation() error {
	if u.Password == "" || u.Email == "" {
		return fmt.Errorf("error password or email is empty")
	}

	return nil
}
