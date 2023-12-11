package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
	"go.uber.org/zap"
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
	var pgErr *pgconn.PgError

	createdUser := model.User{}

	query := `
		INSERT INTO users(company_id, position_id, email, enc_password, active, admin, name, surname, patronymic)
		VALUES($1,$2,$3, $4, $5, $6, $7, $8, $9)
		RETURNING id, company_id, position_id, email, enc_password, active, admin, name, surname, patronymic,
		created_at, updated_at`

	err := u.db.GetContext(ctx, &createdUser, query, val.CompanyID, val.PositionID, val.Email, val.Password,
		val.IsActive, val.IsAdmin, val.Name, val.Surname, val.Patronymic)

	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.ErrEmailAlreadyExists
		}

		return nil, fmt.Errorf("create user: %w", err)
	}

	return &createdUser, nil
}

func (u *uStorage) CreateAdmin(ctx context.Context, admin model.AdminCreate, companyName string) (*model.User, error) {
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

func (u *uStorage) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	query := `SELECT id, company_id, position_id, email, enc_password, active, admin, name, surname, patronymic, 
       		  created_at, updated_at
			  FROM users WHERE email = $1 AND active = true`

	err := u.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &model.User{}, nil
		}

		return &model.User{}, err
	}

	return &user, nil
}
