// Package errs - пакет содержит ошибки приложения.
package errs

import "errors"

var (
	// Базовые ошибки
	ErrCompanyNotFound  = errors.New("Company not found")
	ErrPositionNotFound = errors.New("Position not found")
	ErrCreatorNotFound  = errors.New("Creator does not exist")
	ErrLessonNotFound   = errors.New("Lesson not found")
	ErrCourseNotFound   = errors.New("Course not found")
	ErrCoursesNotFound  = errors.New("Course not found")
	ErrAdminIDNotFound  = errors.New("Admin not found")

	// -- Общие ошибки приложения --

	// ErrInternal - внутренняя ошибка сервера
	ErrInternal = errors.New("Internal error")
	// ErrNotFound - не найдено
	ErrNotFound = errors.New("Record Not found")
	// ErrNotFound - пользователь не найден
	ErrUserNotFound = errors.New("User not found")
	// ErrNoRows - в результате нет записей
	ErrNoRows = errors.New("Sql: no rows in result set")
	// ErrUnauthorized - пользователь не авторизован
	ErrUnauthorized = errors.New("Unauthorized")
	// ErrBadRequest - неверный запрос
	ErrBadRequest = errors.New("Bad request")
	// ErrNotFirstLogin - не первый вход в систему
	ErrNotFirstLogin = errors.New("Not first login")

	// -- Ошибки пользователя --

	// ErrUserAlreadyExists - email должен быть уникальный
	ErrEmailAlreadyExists = errors.New("Email already exists")
	// ErrUserPositionNotFound - пользователь должен ссылаться на существующую должность
	ErrUserPositionNotFound = ErrPositionNotFound
	// ErrUserCompanyNotFound - пользователь должен ссылаться на существующую компанию
	ErrUserCompanyNotFound = ErrCompanyNotFound

	// -- Ошибки должности --

	// ErrPositionCompanyNotFound - должность должна ссылаться на существующую компанию
	ErrPositionCompanyNotFound = ErrCompanyNotFound
	// ErrPositionsNotFound - позиции не найдены
	ErrPositionsNotFound = errors.New("Positions not found")

	// -- Ошибки курсов --

	// ErrCourseUserNotFound - создатель курса должен быть зарегистрирован в системе
	ErrCourseUserNotFound = ErrCreatorNotFound

	// -- Ошибки уроков --

	// ErrLessonCourseNotFound - урок должен ссылаться на существующий курс
	ErrLessonCourseNotFound = ErrCourseNotFound
	// ErrLessonUserNotFound - создатель урока должен быть зарегистрирован в системе
	ErrLessonUserNotFound = ErrCreatorNotFound

	// ErrTextLessonNotFound - текст должен ссылаться на существующий урок
	ErrTextLessonNotFound = ErrLessonNotFound
	// ErrTextUserNotFound - создатель текста должен быть зарегистрирован в системе
	ErrTextUserNotFound = ErrUserNotFound

	// ErrPictureLessonNotFound - картинка должна ссылаться на существующий урок
	ErrPictureLessonNotFound = ErrLessonNotFound
	// ErrPictureUserNotFound - создатель картинки должен быть зарегистрирован в системе
	ErrPictureUserNotFound = ErrUserNotFound

	// -- Ошибки назначений и прогресса по учебным материалам --

	// ErrPositionCourseUsed - 1 курс на 1 должность можно назначить только один раз
	ErrPositionCourseUsed = errors.New("Course already assigned to position")
	// ErrPositionCourseUsed - 1 курс на 1 сотрудника можно назначить только один раз
	ErrUserCourseUsed = errors.New("Course already assigned to user")

	// ErrAssignNotFound - у пользователь для прогресса по уроку должно быть назначение на соответствующий курс
	ErrAssignNotFound = errors.New("User has no course assign")
	// ErrAssignLessonUsed - прогрессу курса могут соответствовать только уникальные прогрессы урока
	ErrAssignLessonUsed = errors.New("Course progress already has the same lesson progress")
)
