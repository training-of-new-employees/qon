package model

import (
	// "errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Code struct {
	Code string `json:"code"`
}

func (c *Code) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Code, validation.Required, is.Digit, validation.Length(4, 4)),
	)
}
