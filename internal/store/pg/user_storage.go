package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.RepositoryUser = (*uStorage)(nil)

type uStorage struct {
	db *sqlx.DB
	transaction
}

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
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, handleError(err)
	}

	return createdUser, nil
}

// CreateAdmin - создание пользователя (админа).
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

// EditAdmin - меняет данные администратора с заданным ID.
func (u *uStorage) EditAdmin(
	ctx context.Context,
	admin model.AdminEdit,
) (*model.AdminEdit, error) {
	err := u.tx(func(tx *sqlx.Tx) error {
		var companyID int
		query :=
			`UPDATE users	 SET 	name = COALESCE($1, name),
		surname = COALESCE($2, surname),
		patronymic = COALESCE($3, patronymic),
		email = COALESCE($4, email)
	 	WHERE id = $5 RETURNING company_id`
		err := tx.GetContext(ctx, &companyID, query, admin.Name, admin.Surname, admin.Patronymic, admin.Email, admin.ID)
		if err != nil {
			return err
		}
		query = `UPDATE companies SET name = COALESCE($1, name) WHERE id = $2`
		_, err = tx.ExecContext(ctx, query, admin.Company, companyID)
		return err

	})
	return &admin, err

}

func (u *uStorage) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	tx, err := u.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("beginning tx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Log.Warn("err during tx rollback %v", zap.Error(err))
		}
	}()

	query := `SELECT id, company_id, position_id, email, enc_password, active, admin, name, surname, patronymic, 
       		  created_at, updated_at
			  FROM users WHERE id = $1 AND active = true`

	err = u.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}

		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing tx: %w", err)
	}

	return &user, nil
}

func (u *uStorage) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	query := `SELECT id, company_id, position_id, email, enc_password, active, admin, name, surname, patronymic, 
       		  created_at, updated_at
			  FROM users WHERE email = $1`

	err := u.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, handleError(err)
	}

	return &user, nil
}

func (u *uStorage) EditUser(ctx context.Context, edit *model.UserEdit) (*model.UserEdit, error) {
	err := u.tx(func(tx *sqlx.Tx) error {
		query := `UPDATE users SET
	 	name = COALESCE($1, name),
	 	surname = COALESCE($2, surname),
	 	patronymic = COALESCE($3, patronymic),
	 	email = COALESCE($4, email),
	 	position_id = COALESCE($5, position_id),
	 	archived = $6
	 WHERE id = $7`
		_, err := tx.ExecContext(ctx, query, edit.Name, edit.Surname, edit.Patronymic, edit.Email, edit.PositionID, edit.IsArchived, edit.ID)
		return err
	},
	)
	return edit, err
}

// GetCompany - получает информацию о компании по id.
func (u *uStorage) GetCompany(ctx context.Context, id int) (*model.Company, error) {
	var comp *model.Company
	err := u.tx(
		func(tx *sqlx.Tx) error {
			query := `SELECT * FROM companies WHERE id=$1`
			err := tx.GetContext(ctx, comp, query, id)
			return err
		},
	)

	return comp, handleError(err)
}

// SetPasswordAndActivateUser установка пароля и активация пользователя.
func (u *uStorage) SetPasswordAndActivateUser(ctx context.Context, userID int, encPassword string) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return fmt.Errorf("beginning tx: %w", err)
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Log.Warn("err during tx rollback %v", zap.Error(err))
			}
		}
	}()

	// Установка нового пароля
	if err := u.updatePasswordTx(ctx, tx, userID, encPassword); err != nil {
		return err
	}

	// Активация пользователя
	if err := u.activateUserTx(ctx, tx, userID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing tx: %w", err)
	}

	return nil
}

// GetUsersByCompany - получает информацию обо всех пользователях в компании.
func (u *uStorage) GetUsersByCompany(ctx context.Context, companyID int) ([]model.User, error) {
	var users []model.User

	err := u.tx(func(tx *sqlx.Tx) error {
		query := `SELECT * FROM users WHERE company_id = $1`
		return tx.SelectContext(ctx, users, query, companyID)
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *uStorage) UpdateUserPassword(ctx context.Context, userID int, password string) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return fmt.Errorf("beginning tx: %w", err)
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Log.Warn("err during tx rollback %v", zap.Error(err))
			}
		}
	}()

	// Установка нового пароля
	if err := u.updatePasswordTx(ctx, tx, userID, password); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing tx: %w", err)
	}

	return nil
}
