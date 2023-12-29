package pg

import (
	"context"
	"math/rand"
	"time"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

func (suite *storeTestSuite) TestCreateUser() {
	suite.NotNil(suite.store)

	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	uniqueEmail := "some@example.org"

	testCases := []struct {
		name string
		data func() model.UserCreate
		err  error
	}{
		{
			name: "success",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = company.ID
				u.PositionID = position.ID
				u.Email = uniqueEmail

				return u
			},
			err: nil,
		},
		{
			name: "unique email",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = company.ID
				u.PositionID = position.ID
				u.Email = uniqueEmail

				return u
			},
			err: errs.ErrEmailAlreadyExists,
		},
		{
			name: "not reference company id",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = 0

				return u
			},
			err: errs.ErrCompanyReference,
		},
		{
			name: "not reference position id",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = company.ID
				u.PositionID = 0

				return u
			},
			err: errs.ErrPositionReference,
		},
		{
			name: "empty email",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = company.ID
				u.PositionID = position.ID
				u.Email = ""

				return u
			},
			err: errs.ErrEmailNotEmpty,
		},

		{
			name: "empty password",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = company.ID
				u.PositionID = position.ID
				u.Password = ""

				return u
			},
			err: errs.ErrPasswordNotEmpty,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			user := tc.data()
			_, err := suite.store.UserStorage().CreateUser(context.TODO(), user)

			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestCreateAdmin() {
	suite.NotNil(suite.store)

	uniqueEmail := "some@example.org"

	testCases := []struct {
		name string
		data func() (string, model.UserCreate)
		err  error
	}{
		{
			name: "success",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()
				u.Email = uniqueEmail

				return "test-company", u
			},
			err: nil,
		},
		{
			name: "success",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()
				u.Email = uniqueEmail

				return "test-company", u
			},
			err: errs.ErrEmailAlreadyExists,
		},
		{
			name: "empty company name",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()

				return "", u
			},
			err: errs.ErrCompanyNameNotEmpty,
		},
		{
			name: "empty email",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()
				u.Email = ""

				return "test-company", u
			},
			err: errs.ErrEmailNotEmpty,
		},
		{
			name: "empty password",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()
				u.Password = ""

				return "test-company", u
			},
			err: errs.ErrPasswordNotEmpty,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			companyName, user := tc.data()
			_, err := suite.store.UserStorage().CreateAdmin(context.TODO(), user, companyName)

			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestEditUser() {
	suite.NotNil(suite.store)

	// Cоздаём компанию
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	// Cоздаём должность
	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	u := model.NewTestUserCreate()
	u.CompanyID = company.ID
	u.PositionID = position.ID
	// Cоздаём пользователя
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), u)
	suite.NoError(err)
	suite.NotEmpty(user)

	testCases := []struct {
		name string
		data func() (model.UserEdit, *model.User)
		err  error
	}{
		{
			name: "edit only email",
			data: func() (model.UserEdit, *model.User) {
				editField := model.UserEdit{ID: user.ID}

				// Данные для редактирования
				editedEmail := "new@newemail.org"
				editField.Email = &editedEmail

				// Ожидаемые данные после редактирования
				expected := user
				expected.Email = editedEmail

				return editField, expected
			},
			err: nil,
		},
		{
			name: "edit only name patronymic surname",
			data: func() (model.UserEdit, *model.User) {
				editField := model.UserEdit{ID: user.ID}

				// Данные для редактирования
				editedName := "Edited"
				editedPatronymic := "Edited"
				editedSurname := "Edited"

				editField.Name = &editedName
				editField.Patronymic = &editedPatronymic
				editField.Surname = &editedSurname

				// Ожидаемые данные после редактирования
				expected := user
				expected.Name = editedName
				expected.Patronymic = editedName
				expected.Surname = editedName

				return editField, expected
			},
			err: nil,
		},
		{
			name: "edit only archived",
			data: func() (model.UserEdit, *model.User) {
				editField := model.UserEdit{ID: user.ID}

				// Данные для редактирования
				editedArchived := true
				editField.IsArchived = &editedArchived

				// Ожидаемые данные после редактирования
				expected := user
				expected.IsArchived = editedArchived

				return editField, expected
			},
			err: nil,
		},
		{
			name: "edit only active",
			data: func() (model.UserEdit, *model.User) {
				editField := model.UserEdit{ID: user.ID}

				// Данные для редактирования
				editedActive := false
				editField.IsActive = &editedActive

				// Ожидаемые данные после редактирования
				expected := user
				expected.IsActive = editedActive

				return editField, expected
			},

			err: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Данные пользователя до редактирования
			before, err := suite.store.UserStorage().GetUserByID(context.TODO(), user.ID)
			suite.NoError(err)
			suite.NotEmpty(before)

			editUser, expected := tc.data()

			// Редактирование записи пользователя
			_, err = suite.store.UserStorage().EditUser(context.TODO(), &editUser)
			suite.Equal(err, tc.err)

			// Данные пользователя после редактирования
			after, err := suite.store.UserStorage().GetUserByID(context.TODO(), user.ID)
			suite.NoError(err)
			suite.NotEmpty(after)

			suite.NotEqual(*before, *after)
			suite.Equal(*after, *expected)
		})
	}
}

func (suite *storeTestSuite) TestEditAdmin() {
	suite.NotNil(suite.store)

	user, err := suite.store.UserStorage().CreateAdmin(context.TODO(), model.NewTestUserCreate(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(user)

	company, err := suite.store.CompanyStorage().GetCompany(context.TODO(), user.CompanyID)

	testCases := []struct {
		name string
		data func() (model.AdminEdit, *model.User, *model.Company)
		err  error
	}{
		{
			name: "edit only email",
			data: func() (model.AdminEdit, *model.User, *model.Company) {
				editField := model.AdminEdit{ID: user.ID}

				// Данные для редактирования
				editedEmail := "new@newemail.org"
				editField.Email = &editedEmail

				// Ожидаемые данные после редактирования
				expected := user
				expected.Email = editedEmail

				return editField, expected, nil
			},
			err: nil,
		},
		{
			name: "edit only name patronymic surname",
			data: func() (model.AdminEdit, *model.User, *model.Company) {
				editField := model.AdminEdit{ID: user.ID}

				// Данные для редактирования
				editedName := "Edited"
				editedPatronymic := "Edited"
				editedSurname := "Edited"

				editField.Name = &editedName
				editField.Patronymic = &editedPatronymic
				editField.Surname = &editedSurname

				// Ожидаемые данные после редактирования
				expected := user
				expected.Name = editedName
				expected.Patronymic = editedName
				expected.Surname = editedName

				return editField, expected, nil
			},
			err: nil,
		},
		{
			name: "edit only company name",
			data: func() (model.AdminEdit, *model.User, *model.Company) {
				editField := model.AdminEdit{ID: user.ID}

				// Данные для редактирования
				editedCompanyName := "new-company-name"
				editField.Company = &editedCompanyName

				// Ожидаемые данные после редактирования
				expected := company
				expected.Name = editedCompanyName

				return editField, nil, expected
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			editUser, expectedUser, expectedCompany := tc.data()

			// Редактирование записи пользователя
			_, err = suite.store.UserStorage().EditAdmin(context.TODO(), editUser)

			// Данные пользователя после редактирования
			afterUser, err := suite.store.UserStorage().GetUserByID(context.TODO(), user.ID)
			suite.NoError(err)
			suite.NotEmpty(afterUser)

			afterCompany, err := suite.store.CompanyStorage().GetCompany(context.TODO(), user.CompanyID)
			suite.NoError(err)
			suite.NotEmpty(afterCompany)

			if expectedUser != nil {
				// suite.NotEqual(*beforeCompany, *afterCompany)
				suite.Equal(*afterUser, *expectedUser)
			}

			if expectedCompany != nil {
				suite.Equal(*afterCompany, *expectedCompany)
			}
		})
	}
}

func (suite *storeTestSuite) TestGetUserByID() {
	suite.NotNil(suite.store)

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	u := model.NewTestUserCreate()
	u.CompanyID = company.ID
	u.PositionID = position.ID

	user, err := suite.store.UserStorage().CreateUser(context.TODO(), u)
	suite.NoError(err)
	suite.NotEmpty(user)

	testCases := []struct {
		name   string
		userID func() int
		err    error
	}{
		{
			name: "success",
			userID: func() int {
				return user.ID
			},
			err: nil,
		},
		{
			name: "random user id",
			userID: func() int {
				return rnd.Intn(32) + 100
			},
			err: errs.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			userID := tc.userID()
			_, err := suite.store.UserStorage().GetUserByID(context.TODO(), userID)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestGetUserByEmail() {
	suite.NotNil(suite.store)

	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	u := model.NewTestUserCreate()
	u.CompanyID = company.ID
	u.PositionID = position.ID

	user, err := suite.store.UserStorage().CreateUser(context.TODO(), u)
	suite.NoError(err)
	suite.NotEmpty(user)

	testCases := []struct {
		name      string
		userEmail func() string
		err       error
	}{
		{
			name: "success",
			userEmail: func() string {
				return user.Email
			},
			err: nil,
		},
		{
			name: "random user email",
			userEmail: func() string {
				return model.NewTestUserCreate().Email
			},
			err: errs.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			email := tc.userEmail()
			_, err := suite.store.UserStorage().GetUserByEmail(context.TODO(), email)
			suite.Equal(tc.err, err)
		})
	}
}
