package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/training-of-new-employees/qon/internal/errs"
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
	if err := validation.Validate(&c.Email, validation.Required); err != nil {
		return errs.ErrEmailNotEmpty
	}

	if err := validation.Validate(&c.Email, is.Email); err != nil {
		return errs.ErrInvalidEmail
	}

	return nil
}
