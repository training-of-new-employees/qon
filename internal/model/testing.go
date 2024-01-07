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

func NewTestListUsers(companyID int) []User {
	number := randomseq.RandomTestInt()
	users := make([]User, number)

	for i := 0; i < number; i++ {
		users[i] = User{
			ID: i + 1, IsActive: true, IsArchived: false,
			CompanyID: companyID, PositionID: randomseq.RandomTestInt() - 99,
			Email:      fmt.Sprintf("%s@example.org", randomseq.RandomString(10)),
			Name:       randomseq.RandomString(10),
			Surname:    randomseq.RandomString(10),
			Patronymic: randomseq.RandomString(10),
		}
	}

	return users
}

func NewTestUser(userID int, companyID int, positionID int) *UserInfo {
	userInfo := &UserInfo{
		User: User{
			ID: userID, IsActive: true, IsArchived: false,
			CompanyID: companyID, PositionID: positionID,
			Email:      fmt.Sprintf("%s@example.org", randomseq.RandomString(10)),
			Name:       randomseq.RandomString(10),
			Surname:    randomseq.RandomString(10),
			Patronymic: randomseq.RandomString(10),
		},
		CompanyName:  fmt.Sprintf("company-name-%s", randomseq.RandomString(5)),
		PositionName: fmt.Sprintf("company-name-%s", randomseq.RandomString(5)),
	}

	return userInfo
}

func NewTestEditUser(userID int, companyID int, positionID int) (UserEdit, UserEdit) {
	user := NewTestUserCreate()

	expected := UserEdit{
		ID:         userID,
		CompanyID:  &companyID,
		PositionID: &positionID,
		Email:      &user.Email,
		Name:       &user.Name,
		Patronymic: &user.Patronymic,
		Surname:    &user.Surname,
		IsActive:   &user.IsActive,
		IsArchived: &user.IsArchived,
	}

	editField := UserEdit{
		ID:         userID,
		CompanyID:  &companyID,
		PositionID: &positionID,
	}

	// определение случайным образом полей для редактирования:
	//
	// изменение емейла
	if randomseq.RandomBool() {
		email := fmt.Sprintf("%s@example.org", randomseq.RandomString(5))

		editField.Email = &email

		expected.Email = &email
	}
	// изменение имени
	if randomseq.RandomBool() {
		name := randomseq.RandomString(8)

		editField.Name = &name

		expected.Name = &name
	}
	// изменение отчества
	if randomseq.RandomBool() {
		patronymic := randomseq.RandomString(8)

		editField.Patronymic = &patronymic

		expected.Patronymic = &patronymic
	}
	// изменение фамилии
	if randomseq.RandomBool() {
		surname := randomseq.RandomString(8)

		editField.Surname = &surname

		expected.Surname = &surname
	}
	// изменение статуса архив
	if randomseq.RandomBool() {
		archived := randomseq.RandomBool()

		editField.IsArchived = &archived

		expected.IsArchived = &archived
	}
	// изменение статуса активности
	if randomseq.RandomBool() {
		active := randomseq.RandomBool()

		editField.IsActive = &active

		expected.IsActive = &active
	}

	return editField, expected
}
