package model

import (
	"github.com/training-of-new-employees/qon/internal/errs"
)

var (
	ErrEmailAlreadyExists = errs.ErrEmailAlreadyExists
	ErrNoRows             = errs.ErrNoRows
	ErrUserNotFound       = errs.ErrUserNotFound
	ErrCompanyIDNotFound  = errs.ErrCompanyNotFound
	ErrPositionNotFound   = errs.ErrPositionNotFound
	ErrPositionsNotFound  = errs.ErrPositionsNotFound
	ErrNoAuthorized       = errs.ErrUnauthorized
)
