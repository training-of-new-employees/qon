// Package errs - пакет содержит ошибки приложения.
package errs

import "errors"

var (
	// Базовые ошибки
	ErrCompanyNotFound  = errors.New("Company not found")
	ErrPositionNotFound = errors.New("Position not found")

	ErrCreatorNotFound = errors.New("Creator does not exist")

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
	// ErrUserPositionNotExist - пользователь должен ссылаться на существующую должность
	ErrUserPositionNotExist = ErrPositionNotFound
	// ErrUserCompanyNotExist - пользователь должен ссылаться на существующую компанию
	ErrUserCompanyNotExist = ErrCompanyNotFound

	// -- Ошибки должности --

	// ErrCompanyNotExist - должность должна ссылаться на существующую компанию
	ErrPositionCompanyNotExist = ErrCompanyNotFound
	// ErrCompanyNotExist - позиции не найдены
	ErrPositionsNotFound = errors.New("Positions not found")

	// -- Ошибки курса --

	// ErrCourseUserNotExist - создатель курса должен быть зарегистрирован в системе
	ErrCourseUserNotExist = ErrCreatorNotFound

	// -- Ошибки урока --

	// ErrLessonCourse - урока должен ссылаться на существующий курс
	ErrLessonCourseNotExist = errors.New("Course does not exist")
	// ErrLessonUserNotExist - создатель урока должен быть зарегистрирован в системе
	ErrLessonUserNotExist = ErrCreatorNotFound
)
