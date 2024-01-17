package pg

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
	"go.uber.org/zap"
)

// transaction - структура для встраивания в репозитории (companyStorage, positionStorage, uStorage,...).
// ВАЖНО: Самостоятельно не используется. Cодержит методы, используемые только в транзакциях.
type transaction struct {
	db *sqlx.DB
}

// tx - обёртка для простого использования транзакций без дублирования кода.
func (tn *transaction) tx(f func(*sqlx.Tx) error) error {
	// открываем транзакцию
	tx, err := tn.db.Beginx()
	if err != nil {
		return fmt.Errorf("beginning tx: %w", err)
	}
	// отмена транзакции
	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Log.Warn("err during tx rollback %v", zap.Error(err))
		}
	}()

	if err = f(tx); err != nil {
		return err
	}

	// фиксация транзакции
	return tx.Commit()
}

// createUserTx - создание пользователя.
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) createUserTx(ctx context.Context, tx *sqlx.Tx, user model.UserCreate) (*model.User, error) {
	createdUser := model.User{}

	query :=
		`INSERT INTO users (
			company_id, position_id, active, admin,
			email, enc_password,
			name, patronymic, surname)
		VALUES($1,$2,$3, $4, $5, $6, $7, $8, $9)
		RETURNING
			id, company_id, position_id,
			active, admin, 
			email, enc_password, name, patronymic, surname,
			created_at, updated_at`

	err := tx.GetContext(
		ctx,
		&createdUser,
		query,
		user.CompanyID,
		user.PositionID,
		user.IsActive,
		user.IsAdmin,
		user.Email,
		user.Password,
		user.Name,
		user.Patronymic,
		user.Surname,
	)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

// updateUserTx - обновление данных пользователя.
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) updateUserTx(ctx context.Context, tx *sqlx.Tx, edit model.UserEdit) (*model.User, error) {
	user := model.User{}

	query := `
		UPDATE
			users
		SET
			company_id  = COALESCE($1, company_id),
			position_id = COALESCE($2, position_id),
			active 	    = COALESCE($3, active),
			archived 	= COALESCE($4, archived),
			email       = COALESCE($5, email),
			name        = COALESCE($6, name),
			patronymic  = COALESCE($7, patronymic),
			surname     = COALESCE($8, surname)
		WHERE id = $9
		RETURNING
			id, company_id, position_id,
			active, archived, admin, 
			email, enc_password, name, patronymic, surname,
			created_at, updated_at
	`

	err := tx.GetContext(
		ctx, &user, query,
		edit.CompanyID, edit.PositionID,
		edit.IsActive, edit.IsArchived,
		edit.Email, edit.Name, edit.Patronymic, edit.Surname, edit.ID,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// createCompanyTx - создание компании в транзакции.
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) createCompanyTx(ctx context.Context, tx *sqlx.Tx, companyName string) (*model.Company, error) {
	company := model.Company{}

	query := `INSERT INTO companies (name) VALUES ($1) RETURNING id, name`

	if err := tx.GetContext(ctx, &company, query, companyName); err != nil {
		return nil, err
	}

	return &company, nil
}

// createCompanyTx - создание компании в транзакции.
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) updateCompanyTx(ctx context.Context, tx *sqlx.Tx, edit model.CompanyEdit) (*model.Company, error) {
	company := model.Company{}

	query := `
		UPDATE
			companies
		SET
			active = COALESCE($1, active),
			archived = COALESCE($2, archived),
			name = COALESCE($3, name)
		WHERE id = $4
		RETURNING id, active, archived, name
	`

	err := tx.GetContext(
		ctx, &company, query,
		edit.IsActive, edit.IsArchived, edit.Name,
		edit.ID,
	)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

// getCompanyTx - получение компании в транзакции.
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) getCompanyTx(ctx context.Context, tx *sqlx.Tx, id int) (*model.Company, error) {
	company := &model.Company{}

	query := `SELECT * FROM companies WHERE id=$1`
	err := tx.GetContext(ctx, company, query, id)
	if err != nil {
		return nil, err
	}
	return company, err
}

// createPositionTx - создание должности в рамках компании.
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) createPositionTx(ctx context.Context, tx *sqlx.Tx, companyID int, positionName string) (*model.Position, error) {
	position := model.Position{}

	query := `
		INSERT INTO
		positions (company_id, name) VALUES ($1, $2)
		RETURNING id, company_id, active, name, created_at, updated_at
	`

	if err := tx.GetContext(ctx, &position, query, companyID, positionName); err != nil {
		return nil, err
	}

	return &position, nil
}

// updatePasswordTx обновляет пароль пользователя.
// ВAЖНО: использовать только только внутри транзакции.
func (tn *transaction) updatePasswordTx(ctx context.Context, tx *sqlx.Tx, userID int, encPassword string) error {
	query := `UPDATE users SET enc_password = $1 WHERE id = $2 RETURNING id`
	err := tx.GetContext(ctx, new(model.User), query, encPassword, userID)
	if err != nil {
		return err
	}

	return nil
}

// activateUserTx активирует пользователя.
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) activateUserTx(ctx context.Context, tx *sqlx.Tx, userID int) error {
	query := `UPDATE users SET active = true WHERE id = $1 RETURNING id`
	err := tx.GetContext(ctx, new(model.User), query, userID)
	if err != nil {
		return err
	}
	return nil
}

// createCompanyTx - назначение курса на должность.
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) assignCourseTx(ctx context.Context, tx *sqlx.Tx, positionID int, courseID int) error {
	query := `INSERT INTO position_course (position_id, course_id) VALUES ($1, $2) RETURNING id`

	if err := tx.QueryRowxContext(ctx, query, positionID, courseID).Err(); err != nil {
		return err
	}

	return nil
}

// insertTextsTx - добавление новой строки в таблицу texts
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) insertTextsTx(ctx context.Context,
	tx *sqlx.Tx, lessonID int, content string, userId int) (string, error) {
	var contentIns string
	query := `INSERT INTO 
			  texts (lesson_id, created_by, content)
			  VALUES ($1, $2, $3)
			  RETURNING content`

	err := tx.GetContext(ctx, &contentIns, query, lessonID,
		userId, content)
	if err != nil {
		return "", err
	}
	return contentIns, nil
}

// updateTextsTx - обновление таблицы texts
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) updateTextsTx(ctx context.Context,
	tx *sqlx.Tx, lessonID int, content string) (string, error) {
	var contentUpd string
	query := `UPDATE texts
			  SET content    = COALESCE(NULLIF($1, ''), content)
			  WHERE lesson_id = $2
			  RETURNING content`
	err := tx.GetContext(ctx, &contentUpd, query, content, lessonID)
	if err != nil {
		return "", err
	}
	return contentUpd, nil
}

// insertPicturesTx - добавление новой строки в таблицу pictures
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) insertPicturesTx(ctx context.Context,
	tx *sqlx.Tx, lessonID int, urlPicture string, userId int) (string, error) {
	var urlPictureIns string
	query := `INSERT INTO
			  pictures (lesson_id, created_by, url_picture)
			  VALUES ($1, $2, $3)
			  RETURNING url_picture`

	err := tx.GetContext(ctx, &urlPictureIns, query,
		lessonID, userId, urlPicture)
	if err != nil {
		return "", err
	}
	return urlPictureIns, nil
}

// updatePicturesTx - обновление таблицы pictures
// ВAЖНО: использовать только внутри транзакции.
func (tn *transaction) updatePicturesTx(ctx context.Context,
	tx *sqlx.Tx, lessonID int, urlPicture string) (string, error) {
	var urlPictureUpd string
	query := `UPDATE pictures
		      SET url_picture = COALESCE(NULLIF($1, ''), url_picture)
		      WHERE lesson_id = $2
			  RETURNING url_picture`
	err := tx.GetContext(ctx, &urlPictureUpd, query, urlPicture, lessonID)
	if err != nil {
		return "", err
	}
	return urlPictureUpd, nil
}
