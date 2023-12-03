package model

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("err email already exists")
	ErrNoRows             = errors.New("err sql: no rows in result set")
	ErrUserNotFound       = errors.New("not found")
)
