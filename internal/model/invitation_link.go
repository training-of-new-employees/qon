package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	InvitationLinkResponse struct {
		Email string `json:"email"`
		Link  string `json:"link"`
	}

	InvitationLinkRequest struct {
		Email string `json:"email"`
	}
)

func (c *InvitationLinkRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Email, validation.Required, is.Email),
	)
}
