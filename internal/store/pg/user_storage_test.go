package pg

import (
	"context"
	"math/rand"
	"time"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
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

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

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
		data func() (model.UserEdit, *model.User) // получаем данные для редактирования и ожидаемый результат
		err  error
	}{
		{
			name: "not found",
			data: func() (model.UserEdit, *model.User) {
				return model.UserEdit{ID: rnd.Intn(32) + 100}, nil
			},
			err: errs.ErrUserNotFound,
		},
		{
			name: "success",
			data: func() (model.UserEdit, *model.User) {
				editField := model.UserEdit{ID: user.ID}

				// ожидаемые данные пользователя
				expected := *user

				// Определяем случайным образом значения для редактирования данных пользователя:
				// Изменение емейла
				if randomseq.RandomBool() {
					email := "new@newemail.org"

					editField.Email = &email
					expected.Email = email

				}
				// Изменение имени
				if randomseq.RandomBool() {
					name := "somename"

					editField.Name = &name
					expected.Name = name
				}
				// Изменение отчества
				if randomseq.RandomBool() {
					patronymic := "somepatronymic"

					editField.Patronymic = &patronymic
					expected.Patronymic = patronymic
				}
				// Изменение фамилии
				if randomseq.RandomBool() {
					surname := "somesurname"

					editField.Surname = &surname
					expected.Surname = surname
				}
				// Изменение статуса архив
				if randomseq.RandomBool() {
					archived := randomseq.RandomBool()

					editField.IsArchived = &archived
					expected.IsArchived = archived
				}
				// Изменение статуса активности
				if randomseq.RandomBool() {
					active := randomseq.RandomBool()

					editField.IsActive = &active
					expected.IsActive = active
				}

				return editField, &expected
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			editUser, expectedResult := tc.data()

			// Редактирование записи пользователя
			_, err := suite.store.UserStorage().EditUser(context.TODO(), &editUser)
			suite.Equal(err, tc.err)

			// Проверка данных пользователя после редактирования
			after, _ := suite.store.UserStorage().GetUserByID(context.TODO(), editUser.ID)
			if after != nil {
				suite.Equal(expectedResult, after)
			}
		})
	}
}

func (suite *storeTestSuite) TestEditAdmin() {
	suite.NotNil(suite.store)

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	// Cоздаём администратора
	user, err := suite.store.UserStorage().CreateAdmin(context.TODO(), model.NewTestUserCreate(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(user)

	// Получаем данные компании
	company, err := suite.store.CompanyStorage().GetCompany(context.TODO(), user.CompanyID)
	suite.NoError(err)
	suite.NotEmpty(company)

	testCases := []struct {
		name string
		data func() (model.AdminEdit, *model.User, *model.Company) // получаем данные для редактирования и ожидаемый результат
		err  error
	}{
		{
			name: "not found",
			data: func() (model.AdminEdit, *model.User, *model.Company) {
				expectedCompany := *company

				adminEdit := model.AdminEdit{ID: rnd.Intn(32) + 100}

				return adminEdit, nil, &expectedCompany
			},
			err: errs.ErrUserNotFound,
		},
		{
			name: "success",
			data: func() (model.AdminEdit, *model.User, *model.Company) {
				expectedCompany := *company
				expectedUser := *user

				adminEdit := model.AdminEdit{ID: user.ID}

				// Определяем случайным образом значения для редактирования данных администратора:
				// Изменение емейла
				if randomseq.RandomBool() {
					email := "new@newemail.org"

					adminEdit.Email = &email
					expectedUser.Email = email

				}
				// Изменение имени
				if randomseq.RandomBool() {
					name := "somename"

					adminEdit.Name = &name
					expectedUser.Name = name
				}
				// Изменение отчества
				if randomseq.RandomBool() {
					patronymic := "somepatronymic"

					adminEdit.Patronymic = &patronymic
					expectedUser.Patronymic = patronymic
				}
				// Изменение фамилии
				if randomseq.RandomBool() {
					surname := "somesurname"

					adminEdit.Surname = &surname
					expectedUser.Surname = surname
				}
				// Изменение имени компании
				if randomseq.RandomBool() {
					companyName := randomseq.RandomString(10)

					adminEdit.Company = &companyName
					expectedCompany.Name = companyName
				}

				return adminEdit, &expectedUser, &expectedCompany
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			adminEdit, expectedUser, expectedCompany := tc.data()

			// Редактирование записи пользователя
			_, err = suite.store.UserStorage().EditAdmin(context.TODO(), adminEdit)
			suite.Equal(tc.err, err)

			// Данные пользователя после редактирования
			userAfter, _ := suite.store.UserStorage().GetUserByID(context.TODO(), adminEdit.ID)
			if userAfter != nil {
				suite.Equal(expectedUser, userAfter)
			}

			// Название компании после редактирования
			companyAfter, err := suite.store.CompanyStorage().GetCompany(context.TODO(), user.CompanyID)
			suite.NoError(err)
			suite.Equal(expectedCompany.Name, companyAfter.Name)

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

func (suite *storeTestSuite) TestGetUsersByCompany() {
	suite.NotNil(suite.store)

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	// Поиск пользователей в пустой базе по id несуществующей компании
	users, err := suite.store.UserStorage().GetUsersByCompany(context.TODO(), rnd.Intn(32)+100)
	suite.Error(err)
	suite.Empty(users)

	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	countUsers := rnd.Intn(32) + 2
	expectedIDs := []int{}

	// Добавление случайного кол-ва пользователей
	for i := 0; i < countUsers; i++ {

		u := model.NewTestUserCreate()
		u.CompanyID = company.ID
		u.PositionID = position.ID

		user, err := suite.store.UserStorage().CreateUser(
			context.TODO(),
			u,
		)
		suite.NoError(err)

		// Добавление в массив идентификаторов добавленных пользователей
		expectedIDs = append(expectedIDs, user.ID)
	}

	// Получение добавленных пользователей
	users, err = suite.store.UserStorage().GetUsersByCompany(context.TODO(), company.ID)
	suite.NotEmpty(users)
	suite.NoError(err)

	// Добавление в массив идентификаторов полученных пользователей
	actualIDs := []int{}
	for _, u := range users {
		actualIDs = append(actualIDs, u.ID)
	}

	suite.EqualValues(expectedIDs, actualIDs)
}
