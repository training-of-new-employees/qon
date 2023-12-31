package model

import (
	"fmt"

	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
)

func NewTestUserCreate() UserCreate {
	return UserCreate{
		Email:      fmt.Sprintf("%s@example.org", randomseq.RandomHexString(16)),
		Password:   randomseq.RandomHexString(64),
		Name:       "Test",
		Patronymic: "Test",
		Surname:    "Test",
		IsAdmin:    false,
		IsActive:   true,
	}
}
