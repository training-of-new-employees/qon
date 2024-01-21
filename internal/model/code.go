package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/training-of-new-employees/qon/internal/errs"
)

type Code struct {
	Code string `json:"code"`
}

func (c *Code) Validate() error {
	// проверка на пустоту поля код верификации
	if err := validation.Validate(&c.Code, validation.Required); err != nil {
		return errs.ErrVerifyCodeNotEmpty
	}

	// проверка на корректность кода верификации
	if err := validation.Validate(&c.Code, is.Digit, validation.Length(4, 4)); err != nil {
		return errs.ErrIncorrectVerifyCode
	}

	return nil
}
