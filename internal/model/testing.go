package model

import (
	"fmt"

	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
)

func NewTestUserCreate() UserCreate {
	return UserCreate{
		Email:      fmt.Sprintf("%s@example.org", randomseq.RandomHexString(16)),
		Password:   randomseq.RandomHexString(64),
		Name:       fmt.Sprintf("Test%s", randomseq.RandomHexString(10)),
		Patronymic: fmt.Sprintf("Test%s", randomseq.RandomHexString(10)),
		Surname:    fmt.Sprintf("Test%s", randomseq.RandomHexString(10)),
		IsAdmin:    false,
		IsActive:   true,
	}
}
