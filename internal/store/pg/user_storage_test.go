package pg

import (
	"context"
	"fmt"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
)

func (suite *storeTestSuite) TestCreateUser() {
	suite.NotNil(suite.store)

	// добавление компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	// добавление должности
	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	// емейл для проверки уникальности поля 'email'
	uniqueEmail := fmt.Sprintf("%s@example.org", randomseq.RandomString(5))

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

				// 2-ой раз используем емейл из предыдущего кейса (кейс "success")
				// для проверки уникальности поля 'email'
				u.Email = uniqueEmail

				return u
			},
			err: errs.ErrEmailAlreadyExists,
		},
		{
			name: "not reference company id",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = randomseq.RandomTestInt()

				return u
			},
			err: errs.ErrCompanyReference,
		},
		{
			name: "not reference position id",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = company.ID
				u.PositionID = randomseq.RandomTestInt()

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

	// емейл для проверки уникальности поля 'email'
	uniqueEmail := fmt.Sprintf("%s@example.org", randomseq.RandomString(5))

	testCases := []struct {
		name string
		data func() (string, model.UserCreate) // возвращает название компании и данные администратора для создания
		err  error
	}{
		{
			name: "success",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()
				u.Email = uniqueEmail

				return "test-company-1", u
			},
			err: nil,
		},
		{
			name: "unique email",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()
				u.Email = uniqueEmail

				return "test-company-2", u
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

				return "test-company-3", u
			},
			err: errs.ErrEmailNotEmpty,
		},
		{
			name: "empty password",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()
				u.Password = ""

				return "test-company-4", u
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

	// добавление компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	// добавление должности
	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	u := model.NewTestUserCreate()
	u.CompanyID = company.ID
	u.PositionID = position.ID
	// добавление пользователя
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), u)
	suite.NoError(err)
	suite.NotEmpty(user)

	testCases := []struct {
		name string
		data func() (model.UserEdit, *model.User) // возвращает данные для редактирования и ожидаемый результат для проверки тест-кейса
		err  error
	}{
		{
			name: "not found",
			data: func() (model.UserEdit, *model.User) {
				return model.UserEdit{ID: randomseq.RandomTestInt()}, nil
			},
			err: errs.ErrUserNotFound,
		},
		{
			name: "empty email",
			data: func() (model.UserEdit, *model.User) {
				editField := model.UserEdit{ID: user.ID}

				// ожидаемые данные пользователя
				expected := *user

				email := ""

				editField.Email = &email
				expected.Email = email

				return editField, &expected
			},
			err: errs.ErrEmailNotEmpty,
		},
		{
			name: "success",
			data: func() (model.UserEdit, *model.User) {
				editField := model.UserEdit{ID: user.ID}

				// ожидаемые данные пользователя
				expected := *user

				// определение случайным образом полей для редактирования:
				//
				// изменение емейла
				if randomseq.RandomBool() {
					email := fmt.Sprintf("%s@example.org", randomseq.RandomString(5))

					editField.Email = &email
					expected.Email = email

				}
				// изменение имени
				if randomseq.RandomBool() {
					name := randomseq.RandomString(8)

					editField.Name = &name
					expected.Name = name
				}
				// изменение отчества
				if randomseq.RandomBool() {
					patronymic := randomseq.RandomString(8)

					editField.Patronymic = &patronymic
					expected.Patronymic = patronymic
				}
				// изменение фамилии
				if randomseq.RandomBool() {
					surname := randomseq.RandomString(8)

					editField.Surname = &surname
					expected.Surname = surname
				}
				// изменение статуса архив
				if randomseq.RandomBool() {
					archived := randomseq.RandomBool()

					editField.IsArchived = &archived
					expected.IsArchived = archived
				}
				// изменение статуса активности
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

			// используется только при успешном кейсе
			if err == nil && after != nil {
				suite.Equal(expectedResult, after)
			}
		})
	}
}

func (suite *storeTestSuite) TestEditAdmin() {
	suite.NotNil(suite.store)

	// добавление администратора
	user, err := suite.store.UserStorage().CreateAdmin(context.TODO(), model.NewTestUserCreate(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(user)

	// получение данных компании
	company, err := suite.store.CompanyStorage().GetCompany(context.TODO(), user.CompanyID)
	suite.NoError(err)
	suite.NotEmpty(company)

	testCases := []struct {
		name string
		data func() (model.AdminEdit, *model.User, *model.Company) // получаем данные для редактирования и ожидаемый результат для проверки тест-кейса
		err  error
	}{
		{
			name: "not found",
			data: func() (model.AdminEdit, *model.User, *model.Company) {
				expectedCompany := *company

				adminEdit := model.AdminEdit{ID: randomseq.RandomTestInt()}

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

				// определение случайным образом полей для редактирования:
				//
				// изменение емейла
				if randomseq.RandomBool() {
					email := fmt.Sprintf("%s@example.org", randomseq.RandomString(5))

					adminEdit.Email = &email
					expectedUser.Email = email

				}
				// изменение имени
				if randomseq.RandomBool() {
					name := randomseq.RandomString(8)

					adminEdit.Name = &name
					expectedUser.Name = name
				}
				// изменение отчества
				if randomseq.RandomBool() {
					patronymic := randomseq.RandomString(8)

					adminEdit.Patronymic = &patronymic
					expectedUser.Patronymic = patronymic
				}
				// изменение фамилии
				if randomseq.RandomBool() {
					surname := randomseq.RandomString(8)

					adminEdit.Surname = &surname
					expectedUser.Surname = surname
				}
				// изменение названия компании
				if randomseq.RandomBool() {
					companyName := randomseq.RandomString(8)

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

			// редактирование записи пользователя
			_, err = suite.store.UserStorage().EditAdmin(context.TODO(), adminEdit)
			suite.Equal(tc.err, err)

			// данные пользователя после редактирования
			userAfter, _ := suite.store.UserStorage().GetUserByID(context.TODO(), adminEdit.ID)
			if userAfter != nil {
				suite.Equal(expectedUser, userAfter)
			}

			// название компании после редактирования
			companyAfter, err := suite.store.CompanyStorage().GetCompany(context.TODO(), user.CompanyID)
			suite.NoError(err)
			suite.Equal(expectedCompany.Name, companyAfter.Name)

		})
	}
}

func (suite *storeTestSuite) TestGetUserByID() {
	suite.NotNil(suite.store)

	// добавление компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	// добавление должности
	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	// подготовка данных пользователя для добавления
	u := model.NewTestUserCreate()
	u.CompanyID = company.ID
	u.PositionID = position.ID

	// добавление пользователя
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), u)
	suite.NoError(err)
	suite.NotEmpty(user)

	testCases := []struct {
		name   string
		userID int
		err    error
	}{
		{
			name:   "success",
			userID: user.ID,
			err:    nil,
		},
		{
			name:   "random user id",
			userID: randomseq.RandomTestInt(),
			err:    errs.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.store.UserStorage().GetUserByID(context.TODO(), tc.userID)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestGetUserByEmail() {
	suite.NotNil(suite.store)

	// добавление компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	// добавление должности
	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	// подготовка данных пользователя для добавления
	u := model.NewTestUserCreate()
	u.CompanyID = company.ID
	u.PositionID = position.ID

	// добавление пользователя
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

	// поиск пользователей в пустой базе по id несуществующей компании
	users, err := suite.store.UserStorage().GetUsersByCompany(context.TODO(), randomseq.RandomTestInt())
	suite.Error(err)
	suite.Empty(users)

	// добавление компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	// добавление должности
	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	// генерация кол-ва пользователей
	countUsers := randomseq.RandomTestInt()

	// массив с ожидаемыми идентификаторами пользователей
	expectedIDs := []int{}

	// добавление случайного кол-ва пользователей (от 100 до 356)
	for i := 0; i < countUsers; i++ {

		u := model.NewTestUserCreate()
		u.CompanyID = company.ID
		u.PositionID = position.ID

		user, err := suite.store.UserStorage().CreateUser(
			context.TODO(),
			u,
		)
		suite.NoError(err)

		// добавление в массив идентификаторов добавленных пользователей
		expectedIDs = append(expectedIDs, user.ID)
	}

	// получение добавленных пользователей
	users, err = suite.store.UserStorage().GetUsersByCompany(context.TODO(), company.ID)
	suite.NotEmpty(users)
	suite.NoError(err)

	// добавление в массив идентификаторов полученных пользователей
	actualIDs := []int{}
	for _, u := range users {
		actualIDs = append(actualIDs, u.ID)
	}

	suite.EqualValues(expectedIDs, actualIDs)
}

func (suite *storeTestSuite) TestSetPassword() {
	suite.NotNil(suite.store)

	// добавление компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	// добавление должности
	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	// подготовка данных пользователя для добавления
	u := model.NewTestUserCreate()
	u.CompanyID = company.ID
	u.PositionID = position.ID

	// создание пользователя
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), u)
	suite.NoError(err)
	suite.NotEmpty(user)

	testCases := []struct {
		name     string
		userID   int
		password string
		err      error
	}{
		{
			name:     "success",
			userID:   user.ID,
			password: randomseq.RandomString(20),
			err:      nil,
		},
		{
			name:     "not found",
			userID:   randomseq.RandomTestInt(),
			password: randomseq.RandomString(20),
			err:      errs.ErrUserNotFound,
		},
		{
			name:     "empty password",
			userID:   user.ID,
			password: "",
			err:      errs.ErrPasswordNotEmpty,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := suite.store.UserStorage().UpdateUserPassword(context.TODO(), tc.userID, tc.password)
			suite.Equal(tc.err, err)

			userAfter, _ := suite.store.UserStorage().GetUserByID(context.TODO(), tc.userID)
			// используется только в успешном кейсе
			if err == nil && userAfter != nil {
				suite.Equal(tc.password, userAfter.Password)
			}
		})
	}
}

func (suite *storeTestSuite) TestSetPasswordAndActivateUser() {
	suite.NotNil(suite.store)

	// создаём компанию
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	// создаём должность
	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	// подготовка данных пользователя для добавления
	u := model.NewTestUserCreate()
	u.CompanyID = company.ID
	u.PositionID = position.ID
	u.IsActive = false

	// создаём пользователя
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), u)
	suite.NoError(err)
	suite.NotEmpty(user)

	testCases := []struct {
		name     string
		userID   int
		password string
		err      error
	}{
		{
			name:     "success",
			userID:   user.ID,
			password: randomseq.RandomString(20),
			err:      nil,
		},
		{
			name:     "not found",
			userID:   randomseq.RandomTestInt(),
			password: randomseq.RandomString(20),
			err:      errs.ErrUserNotFound,
		},
		{
			name:     "empty password",
			userID:   user.ID,
			password: "",
			err:      errs.ErrPasswordNotEmpty,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := suite.store.UserStorage().SetPasswordAndActivateUser(context.TODO(), tc.userID, tc.password)
			suite.Equal(tc.err, err)

			userAfter, _ := suite.store.UserStorage().GetUserByID(context.TODO(), tc.userID)
			// используется только при успешном кейсе
			if err == nil && userAfter != nil {
				suite.Equal(tc.password, userAfter.Password)
				suite.Equal(true, userAfter.IsActive)
			}
		})
	}
}
