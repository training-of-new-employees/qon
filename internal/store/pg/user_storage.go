package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.RepositoryUser = (*uStorage)(nil)

// uStorage - репозиторий для пользователя (админ, сотрудник).
type uStorage struct {
	db *sqlx.DB
	transaction
}

// newUStorages - конструктор репозитория пользователя.
func newUStorages(db *sqlx.DB) *uStorage {
	return &uStorage{
		db:          db,
		transaction: transaction{db: db},
	}
}

// CreateUser - создание пользователя (сотрудника).
func (u *uStorage) CreateUser(ctx context.Context, val model.UserCreate) (*model.User, error) {
	var createdUser *model.User

	// открываем транзакцию
	err := u.tx(func(tx *sqlx.Tx) error {
		var err error

		// создание пользователя
		createdUser, err = u.createUserTx(ctx, tx, val)

		return err
	})
	if err != nil {
		return nil, handleError(err)
	}

	return createdUser, nil
}

// CreateAdmin - создание администратора.
func (u *uStorage) CreateAdmin(ctx context.Context, admin model.UserCreate, companyName string) (*model.User, error) {
	var createdAdmin *model.User

	// открываем транзакцию
	err := u.tx(func(tx *sqlx.Tx) error {
		var err error

		// создание компании
		company, err := u.createCompanyTx(ctx, tx, companyName)
		if err != nil {
			return err
		}

		// создание должности-заглушки для администратора
		position, err := u.createPositionTx(ctx, tx, company.ID, "admin")
		if err != nil {
			return err
		}

		admin.CompanyID = company.ID
		admin.PositionID = position.ID
		admin.IsAdmin = true

		// создание админа
		createdAdmin, err = u.createUserTx(ctx, tx, admin)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, handleError(err)
	}

	return createdAdmin, nil
}

// GetUserByID - получение данных пользователя по ID.
func (u *uStorage) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User

	query := `
		SELECT
			id, company_id, position_id,
			active, archived, admin,
			email, enc_password, name, patronymic, surname,
			created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := u.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}

		return nil, handleError(err)
	}

	return &user, nil
}

// GetUserByEmail - получение данных пользователя по емейлу.
func (u *uStorage) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	query := `
		SELECT
			id, company_id, position_id,
			active, archived, admin,
			email, enc_password, name, patronymic, surname,
			created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := u.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}

		return nil, handleError(err)
	}

	return &user, nil
}

// GetUsersByCompany - получение данных для каждого пользователя в компании.
func (u *uStorage) GetUsersByCompany(ctx context.Context, companyID int) ([]model.User, error) {
	users := make([]model.User, 0, 10)

	// открываем транзакцию
	err := u.tx(func(tx *sqlx.Tx) error {
		query := `SELECT * FROM users WHERE company_id = $1`
		return tx.SelectContext(ctx, &users, query, companyID)
	})
	if err != nil {
		return nil, handleError(err)
	}
	if len(users) == 0 {
		return nil, errs.ErrUserNotFound
	}

	return users, nil
}

// EditUser - редактирование данных пользователя.
func (u *uStorage) EditUser(ctx context.Context, edit *model.UserEdit) (*model.UserEdit, error) {
	var user *model.User

	// открываем транзакцию
	err := u.tx(func(tx *sqlx.Tx) error {
		var err error

		user, err = u.updateUserTx(ctx, tx, *edit)

		return err
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, handleError(err)
	}

	// TODO: позже нужно исправить, пока используем такое преобразование для совместимости с верхним уровнем
	edit.ID = user.ID
	edit.CompanyID = &user.CompanyID
	edit.PositionID = &user.PositionID
	edit.IsActive = &user.IsActive
	edit.IsArchived = &user.IsArchived
	edit.Email = &user.Email
	edit.Name = &user.Name
	edit.Patronymic = &user.Patronymic
	edit.Surname = &user.Surname

	return edit, nil
}

// EditAdmin - изменение данных администратора.
func (u *uStorage) EditAdmin(ctx context.Context, admin model.AdminEdit) (*model.AdminEdit, error) {
	var user *model.User
	var company *model.Company

	// открываем транзакцию
	err := u.tx(func(tx *sqlx.Tx) error {
		var err error

		// Изменение данных админа
		user, err = u.updateUserTx(
			ctx, tx,
			model.UserEdit{
				ID: admin.ID, Email: admin.Email,
				Name: admin.Name, Patronymic: admin.Patronymic, Surname: admin.Surname,
			},
		)
		if err != nil {
			return err
		}

		// Изменение названия компании
		company, err = u.updateCompanyTx(
			ctx, tx,
			model.CompanyEdit{ID: user.CompanyID, Name: admin.Company},
		)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}

		return nil, handleError(err)
	}

	// TODO: позже нужно исправить, пока используем такое преобразование для совместимости с верхним уровнем
	admin.ID = user.ID
	admin.Email = &user.Email
	admin.Name = &user.Name
	admin.Patronymic = &user.Patronymic
	admin.Surname = &user.Surname
	admin.Company = &company.Name

	return &admin, nil
}

// UpdateUserPassword - обновление пароля пользователя.
func (u *uStorage) UpdateUserPassword(ctx context.Context, userID int, password string) error {
	err := u.tx(func(tx *sqlx.Tx) error {
		// Установка нового пароля
		return u.updatePasswordTx(ctx, tx, userID, password)
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrUserNotFound
		}
		return handleError(err)
	}

	return nil
}

// SetPasswordAndActivateUser - установка пароля и активация пользователя.
func (u *uStorage) SetPasswordAndActivateUser(ctx context.Context, userID int, encPassword string) error {
	// открываем транзакцию
	err := u.tx(func(tx *sqlx.Tx) error {
		// Установка нового пароля
		if err := u.updatePasswordTx(ctx, tx, userID, encPassword); err != nil {
			return err
		}

		// Активация пользователя
		if err := u.activateUserTx(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrUserNotFound
		}
		return handleError(err)
	}

	return nil
}
