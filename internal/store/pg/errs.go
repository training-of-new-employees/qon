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
	"fk_position_company": errs.ErrPositionCompanyNotFound,
	"fk_user_company":     errs.ErrUserCompanyNotFound,
	"fk_user_position":    errs.ErrUserPositionNotFound,
	"unq_user_email":      errs.ErrEmailAlreadyExists,
	"fk_course_user":      errs.ErrCourseUserNotFound,
	"fk_lesson_course":    errs.ErrLessonCourseNotFound,
	"fk_lesson_user":      errs.ErrLessonUserNotFound,
	"fk_text_lesson":      errs.ErrTextLessonNotFound,
	"fk_text_user":        errs.ErrTextUserNotFound,
	"fk_picture_lesson":   errs.ErrPictureLessonNotFound,
	"fk_picture_user":     errs.ErrPictureUserNotFound,

	"fk_positioncourse_position": errs.ErrPositionNotFound,
	"fk_positioncourse_course":   errs.ErrCourseNotFound,
	"unq_positioncourse":         errs.ErrPositionCourseUsed,

	"fk_courseassign_user":   errs.ErrUserNotFound,
	"fk_courseassign_course": errs.ErrCourseNotFound,
	"unq_usercourse":         errs.ErrUserCourseUsed,

	"fk_lessonresult_courseassign": errs.ErrAssignNotFound,
	"fk_lessonresult_lesson":       errs.ErrLessonNotFound,
	"unq_assignlesson":             errs.ErrAssignLessonUsed,
}

func handleError(err error) error {
	// Если ошибки нет, возвращаем nil
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
