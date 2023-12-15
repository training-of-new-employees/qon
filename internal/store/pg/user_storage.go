package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.RepositoryUser = (*uStorage)(nil)

type uStorage struct {
	db *sqlx.DB
}

func newUStorages(db *sqlx.DB) *uStorage {
	return &uStorage{
		db: db,
	}
}

func (u *uStorage) CreateUser(ctx context.Context, val model.UserCreate) (*model.User, error) {
	createdUser := model.User{}

	query := `
		INSERT INTO users(company_id, position_id, email, enc_password, active, admin, name, surname, patronymic)
		VALUES($1,$2,$3, $4, $5, $6, $7, $8, $9)
		RETURNING id, company_id, position_id, email, enc_password, active, admin, name, surname, patronymic,
		created_at, updated_at`

	err := u.db.GetContext(
		ctx,
		&createdUser,
		query,
		val.CompanyID,
		val.PositionID,
		val.Email,
		val.Password,
		val.IsActive,
		val.IsAdmin,
		val.Name,
		val.Surname,
		val.Patronymic,
	)
	if err != nil {
		return nil, handleError(err)
	}

	return &createdUser, nil
}

func (u *uStorage) CreateAdmin(
	ctx context.Context,
	admin model.AdminCreate,
	companyName string,
) (*model.User, error) {
	tx, err := u.db.Beginx()
	if err != nil {
		return &model.User{}, fmt.Errorf("beginning tx: %w", err)
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Log.Warn("err during tx rollback %v", zap.Error(err))
			}
		}
	}()

	company := model.Company{}

	query := `INSERT INTO companies(name) VALUES ($1) RETURNING id, name`

	if err = tx.GetContext(ctx, &company, query, companyName); err != nil {
		return &model.User{}, err
	}

	position := model.Position{}
	query = `INSERT INTO positions(name, company_id) VALUES ($1, $2) RETURNING id, name`

	if err = tx.GetContext(ctx, &position, query, "admin", company.ID); err != nil {
		return &model.User{}, err
	}

	createdAdmin := model.User{}

	query = `INSERT INTO users (company_id, position_id, email, enc_password, active, admin, name, surname, patronymic)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			  RETURNING id, company_id, position_id, email, enc_password, active, admin, name, surname, patronymic,
			  created_at, updated_at`

	if err = tx.GetContext(ctx, &createdAdmin, query, company.ID, position.ID, admin.Email, admin.Password,
		admin.IsActive, admin.IsAdmin, admin.Name, admin.Surname, admin.Patronymic); err != nil {

		return &model.User{}, err
	}

	if err = tx.Commit(); err != nil {
		return &model.User{}, fmt.Errorf("committing tx: %w", err)
	}

	return &createdAdmin, nil
}

// EditAdmin - меняет данные администратора с заданным ID
func (u *uStorage) EditAdmin(
	ctx context.Context,
	admin *model.AdminEdit,
) (*model.AdminEdit, error) {
	tx, err := u.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("beginning tx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Log.Warn("err during tx rollback %v", zap.Error(err))
		}
	}()

	var pgErr *pgconn.PgError
	var companyID int

	query :=
		`UPDATE users
	 SET 
	 	name = COALESCE($1, name),
	 	surname = COALESCE($2, surname),
	 	patronymic = COALESCE($3, patronymic),
	 	email = COALESCE($4, email)
	 WHERE id = $5
	 RETURNING company_id`

	if err = tx.GetContext(ctx, &companyID, query, admin.Name, admin.Surname, admin.Patronymic, admin.Email, admin.ID); err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NoDataFound {
			return nil, model.ErrUserNotFound
		}
		return nil, err

	}

	query = `UPDATE companies SET name = COALESCE($1, name) WHERE id = $2`

	if _, err = tx.ExecContext(ctx, query, admin.Company, companyID); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing tx: %w", err)
	}

	return admin, nil

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
		if errors.Is(err, sql.ErrNoRows) {
			return &model.User{}, nil
		}

		return &model.User{}, err
	}

	return &user, nil
}

// GetCompany - получает информацию о компании по id
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

// GetUsersByCompany - получает информацию обо всех пользователях в компании
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

// updatePasswordTx обновляет пароль пользователя.
// ВAЖНО: может вызываться только внутри транзакции.
func (u *uStorage) updatePasswordTx(ctx context.Context, tx *sqlx.Tx, userID int, encPassword string) error {
	query := `UPDATE users SET enc_password = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, encPassword, userID)
	if err != nil {
		return err
	}

	return nil
}

// activateUserTx активирует пользователя.
// ВAЖНО: может вызываться только внутри транзакции.
func (u *uStorage) activateUserTx(ctx context.Context, tx *sqlx.Tx, userID int) error {
	query := `UPDATE users SET active = true WHERE id = $1`
	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

// tx - обёртка для простого использования транзакций без дублирования кода
func (u *uStorage) tx(f func(*sqlx.Tx) error) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return fmt.Errorf("beginning tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Log.Warn("err during tx rollback %v", zap.Error(err))
		}
	}()

	return handleError(f(tx))
}
