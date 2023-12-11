package model

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNoRows             = errors.New("sql: no rows in result set")
	ErrUserNotFound       = errors.New("not found")
	ErrCompanyIDNotFound  = errors.New("company not found")
	ErrPositionNotFound   = errors.New("position not found")
	ErrPositionsNotFound  = errors.New("positions not found")
)
