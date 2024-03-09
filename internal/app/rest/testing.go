package rest

import (
	"fmt"

	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
)

func NewTestReqCreateUser(CompanyID, PositionID int) (reqCreateUser, model.User) {
	request := reqCreateUser{
		Email:      randomseq.RandomEmail(),
		Name:       fmt.Sprintf("Test%s", randomseq.RandomName(10)),
		Patronymic: fmt.Sprintf("Test%s", randomseq.RandomName(10)),
		Surname:    fmt.Sprintf("Test%s", randomseq.RandomName(10)),
		PositionID: PositionID,
	}

	user := model.User{
		ID:         2,
		CompanyID:  CompanyID,
		PositionID: PositionID,
		Email:      request.Email,
		Name:       request.Name,
		Patronymic: request.Patronymic,
		Surname:    request.Surname,
	}

	return request, user
}

func NewTestReqEditUser(userID, companyID, positionID int) (reqEditUser, model.UserEdit) {
	request := reqEditUser{}

	user := model.NewTestUser(userID, companyID, positionID)

	expected := model.UserEdit{
		ID:         userID,
		CompanyID:  &companyID,
		PositionID: &positionID,
		Email:      &user.Email,
		Name:       &user.Name,
		Surname:    &user.Surname,
		IsActive:   &user.IsActive,
		IsArchived: &user.IsArchived,
	}

	// определение случайным образом полей для редактирования:
	//
	// изменение емейла
	if randomseq.RandomBool() {
		email := randomseq.RandomEmail()

		request.Email = &email
		expected.Email = &email
	}
	// изменение имени
	if randomseq.RandomBool() {
		name := randomseq.RandomName(8)

		request.Name = &name
		expected.Name = &name
	}
	// изменение отчества
	if randomseq.RandomBool() {
		patronymic := randomseq.RandomName(8)

		request.Patronymic = &patronymic
		expected.Patronymic = &patronymic
	}
	// изменение фамилии
	if randomseq.RandomBool() {
		surname := randomseq.RandomName(8)

		request.Surname = &surname
		expected.Surname = &surname
	}
	// изменение статуса архив
	if randomseq.RandomBool() {
		archived := randomseq.RandomBool()

		request.IsArchived = &archived
		expected.IsArchived = &archived
	}
	// изменение статуса активности
	if randomseq.RandomBool() {
		active := randomseq.RandomBool()

		request.IsActive = &active
		expected.IsActive = &active
	}

	return request, expected

}
