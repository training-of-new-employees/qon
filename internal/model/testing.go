package model

import (
	"fmt"

	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
)

func NewTestUserCreate() UserCreate {
	hash, _ := GenerateHash(randomseq.RandomPassword())

	return UserCreate{
		Email:      fmt.Sprintf("%s@example.org", randomseq.RandomString(16)),
		Password:   hash,
		Name:       fmt.Sprintf("Test%s", randomseq.RandomName(10)),
		Patronymic: fmt.Sprintf("Test%s", randomseq.RandomName(10)),
		Surname:    fmt.Sprintf("Test%s", randomseq.RandomName(10)),
		IsAdmin:    false,
		IsActive:   true,
	}
}

func NewTestCreateAdmin() CreateAdmin {
	return CreateAdmin{
		Email:    fmt.Sprintf("%s@example.org", randomseq.RandomString(10)),
		Password: randomseq.RandomPassword(),
		Company:  fmt.Sprintf("test-%s", randomseq.RandomName(10)),
	}
}

func NewTestPositionSet() PositionSet {
	return PositionSet{
		Name: fmt.Sprintf("test-%s", randomseq.RandomName(10)),
	}
}

func NewTestPositions(companyID int) []*Position {
	number := randomseq.RandomTestInt()
	positions := make([]*Position, number)

	for i := 0; i < number; i++ {
		positions[i] = &Position{ID: i + 1, CompanyID: companyID, IsActive: true, IsArchived: false, Name: fmt.Sprintf("test-%s", randomseq.RandomName(10))}
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
			Name:       randomseq.RandomName(10),
			Surname:    randomseq.RandomName(10),
			Patronymic: randomseq.RandomName(10),
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
			Name:       randomseq.RandomName(10),
			Surname:    randomseq.RandomName(10),
			Patronymic: randomseq.RandomName(10),
		},
		CompanyName:  fmt.Sprintf("company-name-%s", randomseq.RandomName(5)),
		PositionName: fmt.Sprintf("company-name-%s", randomseq.RandomName(5)),
	}

	return userInfo
}

func NewTestEditUser(userID int, companyID int, positionID int) (UserEdit, UserEdit) {
	user := NewTestUser(userID, companyID, positionID)

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
		name := randomseq.RandomName(8)

		editField.Name = &name

		expected.Name = &name
	}
	// изменение отчества
	if randomseq.RandomBool() {
		patronymic := randomseq.RandomName(8)

		editField.Patronymic = &patronymic

		expected.Patronymic = &patronymic
	}
	// изменение фамилии
	if randomseq.RandomBool() {
		surname := randomseq.RandomName(8)

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

func NewTestAdminEdit(userID int, companyID int, positionID int) (AdminEdit, AdminEdit) {
	user := NewTestUser(userID, companyID, positionID)

	expected := AdminEdit{
		ID:         userID,
		Email:      &user.Email,
		Company:    &user.CompanyName,
		Name:       &user.Name,
		Patronymic: &user.Patronymic,
		Surname:    &user.Surname,
	}

	editField := AdminEdit{ID: userID}

	// определение случайным образом полей для редактирования:
	//
	// изменение емейла
	if randomseq.RandomBool() {
		email := fmt.Sprintf("%s@example.org", randomseq.RandomString(5))

		editField.Email = &email
		expected.Email = &email

	}
	// изменение названия компании
	if randomseq.RandomBool() {
		companyName := randomseq.RandomName(8)

		editField.Company = &companyName
		expected.Company = &companyName
	}
	// изменение имени
	if randomseq.RandomBool() {
		name := randomseq.RandomName(8)

		editField.Name = &name
		expected.Name = &name
	}
	// изменение отчества
	if randomseq.RandomBool() {
		patronymic := randomseq.RandomName(8)

		editField.Patronymic = &patronymic
		expected.Patronymic = &patronymic
	}
	// изменение фамилии
	if randomseq.RandomBool() {
		surname := randomseq.RandomName(8)

		editField.Surname = &surname
		expected.Surname = &surname
	}

	return editField, expected
}

func NewTestResetPassword() EmailReset {
	return EmailReset{
		Email: fmt.Sprintf("%s@example.org", randomseq.RandomString(16)),
	}
}

func NewTestCourseSet() CourseSet {
	return CourseSet{
		ID:          randomseq.RandomTestInt(),
		Name:        randomseq.RandomName(minNameL),
		Description: randomseq.RandomName(minDescL),
	}
}
