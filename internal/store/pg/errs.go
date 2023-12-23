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
	// companies
	"chck_company_name_not_empty": errs.ErrCompanyNameNotEmpty,

	// positions
	"chck_position_company_not_empty": errs.ErrCompanyIDNotEmpty,
	"fk_position_company":             errs.ErrCompanyReference,
	"chck_position_name_not_empty":    errs.ErrPositionNameNotEmpty,

	// users
	"chck_user_company_not_empty":     errs.ErrCompanyIDNotEmpty,
	"fk_user_company":                 errs.ErrCompanyReference,
	"chck_user_position_not_empty":    errs.ErrPositionIDNotEmpty,
	"fk_user_position":                errs.ErrPositionReference,
	"chck_user_email_not_empty":       errs.ErrEmailNotEmpty,
	"unq_user_email":                  errs.ErrEmailAlreadyExists,
	"chck_user_encpassword_not_empty": errs.ErrPasswordNotEmpty,

	// courses
	"chck_course_creater_not_empty": errs.ErrCreaterNotEmpty,
	"fk_course_creater":             errs.ErrCreaterNotFound,
	"chck_course_name_not_empty":    errs.ErrCourseNameNotEmpty,

	// lessons
	"chck_lesson_course_not_empty":  errs.ErrCourseIDNotEmpty,
	"fk_lesson_course":              errs.ErrCourseReference,
	"chck_lesson_creater_not_empty": errs.ErrCreaterNotEmpty,
	"fk_lesson_creater":             errs.ErrCreaterNotFound,
	"chck_lesson_name_not_empty":    errs.ErrLessonNameNotEmpty,

	// texts
	"chck_text_lesson_not_empty":  errs.ErrLessonIDNotEmpty,
	"fk_text_lesson":              errs.ErrLessonReference,
	"chck_text_creater_not_empty": errs.ErrCreaterNotEmpty,
	"fk_text_creater":             errs.ErrCreaterNotFound,
	"chck_text_content_not_empty": errs.ErrTextContentNotEmpty,

	// pictures
	"chck_picture_lesson_not_empty":  errs.ErrLessonIDNotEmpty,
	"fk_picture_lesson":              errs.ErrLessonReference,
	"chck_picture_creater_not_empty": errs.ErrCreaterNotEmpty,
	"fk_picture_creater":             errs.ErrCreaterNotFound,
	"chck_link_not_empty":            errs.ErrPictureLinkNotEmpty,

	// position_course
	"chck_positioncourse_position_not_empty": errs.ErrPositionIDNotEmpty,
	"fk_positioncourse_position":             errs.ErrPositionReference,
	"chck_positioncourse_course_not_empty":   errs.ErrCourseIDNotEmpty,
	"fk_positioncourse_course":               errs.ErrCourseReference,
	"unq_positioncourse":                     errs.ErrPositionCourseUsed,

	// course_assign
	"chck_courseassign_course_not_empty": errs.ErrCourseIDNotEmpty,
	"fk_courseassign_course":             errs.ErrCourseReference,
	"chck_courseassign_user_not_empty":   errs.ErrUserIDNotEmpty,
	"fk_courseassign_user":               errs.ErrUserReference,
	"unq_usercourse":                     errs.ErrUserCourseUsed,

	// lesson_results
	"chck_lessonresult_course_not_empty": errs.ErrCourseIDNotEmpty,
	"fk_lessonresult_course":             errs.ErrCourseReference,
	"chck_lessonresult_lesson_not_empty": errs.ErrLessonIDNotEmpty,
	"fk_lessonresult_lesson":             errs.ErrLessonReference,
	"chck_lessonresult_user_not_empty":   errs.ErrUserIDNotEmpty,
	"fk_lessonresult_user":               errs.ErrUserReference,

	"unq_assignlesson": errs.ErrAssignLessonUsed,
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
