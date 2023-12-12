package pg

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/logger"
	"go.uber.org/zap"
)

// constraintToAppError - соответствие ограничений CУБД ошибкам приложения
var constraintToAppError = map[string]error{
	"fk_position_company": errs.ErrPositionCompanyNotExist,
	"fk_user_company":     errs.ErrUserCompanyNotExist,
	"fk_user_position":    errs.ErrUserPositionNotExist,
	"unq_user_email":      errs.ErrEmailAlreadyExists,
	"fk_course_user":      errs.ErrCourseUserNotExist,
	"fk_lesson_course":    errs.ErrLessonCourseNotExist,
	"fk_lesson_user":      errs.ErrLessonUserNotExist,
}

func handleError(err error) error {
	// Если ошибки нет, то возвращаем nil
	if err == nil {
		return nil
	}

	// Если не найдены записи, то возвращаем ErrNotFound
	if errors.Is(err, sql.ErrNoRows) {
		return errs.ErrNotFound
	}

	// Проверка, является ли ошибка нарушением ограничения СУБД
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.ConstraintName != "" {
		// Если ограничение известно, то возвращаем соответствующую ошибку
		if appErr, ok := constraintToAppError[pgErr.ConstraintName]; ok {
			return appErr
		}
	}

	// Если другая ошибка или неизвестное ограничение, то возвращаем ErrInternal
	logger.Log.Warn("internal error: %v", zap.Error(err))
	return errs.ErrInternal
}
