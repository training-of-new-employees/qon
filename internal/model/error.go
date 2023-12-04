package model

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("error email already exists")
	ErrNoRows             = errors.New("error sql: no rows in result set")
	ErrUserNotFound       = errors.New("not found")
)
