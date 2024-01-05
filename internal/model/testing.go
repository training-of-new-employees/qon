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

func NewTestCreateAdmin() CreateAdmin {
	return CreateAdmin{
		Email:    fmt.Sprintf("%s@example.org", randomseq.RandomString(10)),
		Password: "abcdA1*",
		Company:  fmt.Sprintf("test-%s", randomseq.RandomString(10)),
	}
}

func NewTestPositionSet() PositionSet {
	return PositionSet{
		Name: fmt.Sprintf("test-%s", randomseq.RandomString(10)),
	}
}

func NewTestPositions(companyID int) []*Position {
	number := randomseq.RandomTestInt()
	positions := make([]*Position, number)

	for i := 0; i < number; i++ {
		positions[i] = &Position{ID: i + 1, CompanyID: companyID, IsActive: true, IsArchived: false, Name: fmt.Sprintf("test-%s", randomseq.RandomString(10))}
	}

	return positions
}

func NewTestUsers(companyID int) []User {
	number := randomseq.RandomTestInt()
	users := make([]User, number)

	for i := 0; i < number; i++ {
		users[i] = User{
			ID: i + 1, IsActive: true, IsArchived: false,
			CompanyID: companyID, PositionID: randomseq.RandomTestInt() - 99,
			Email:      fmt.Sprintf("%s@example.org", randomseq.RandomString(10)),
			Name:       fmt.Sprintf("%s", randomseq.RandomString(10)),
			Surname:    fmt.Sprintf("%s", randomseq.RandomString(10)),
			Patronymic: fmt.Sprintf("%s", randomseq.RandomString(10)),
		}
	}

	return users
}
