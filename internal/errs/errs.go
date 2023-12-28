// Package errs - пакет содержит ошибки приложения.
package errs

import (
	"errors"
)

var (
	// -- Базовые ошибки приложения --

	// ErrInternal - внутренняя ошибка сервера
	ErrInternal = errors.New("internal error")

	// ErrNotSendEmail - ошибка при отправки емейла пользователю
	ErrNotSendEmail = errors.New("can't send email to user")
	// ErrNotFound - не найдено
	ErrNotFound = errors.New("record not found")
	// ErrNotFound - пользователь не найден
	ErrUserNotFound = errors.New("user not found")
	// ErrNoRows - в результате нет записей
	ErrNoRows = errors.New("no rows in result set")
	// ErrUnauthorized - пользователь не авторизован
	ErrUnauthorized = errors.New("unauthorized")
	// ErrBadRequest - неверный запрос
	ErrBadRequest = errors.New("bad request")
	// ErrInvalidRequest - невалидное тело запроса
	ErrInvalidRequest = errors.New("invalid request body")
	// ErrNotFirstLogin - не первый вход в систему
	ErrNotFirstLogin = errors.New("not first login")
	// ErrOnlyAdmin - действие доступно только администратору
	ErrOnlyAdmin = errors.New("you aren't admin")
)

var (
	// -- Общие ошибки объектов --

	// ErrCompanyNotFound - компания не найдена
	ErrCompanyNotFound = errors.New("company not found")
	// ErrPositionNotFound - должность не найдена
	ErrPositionNotFound = errors.New("position not found")

	// ErrLessonNotFound - урок не найден
	ErrLessonNotFound = errors.New("lesson not found")

	// ErrCourseNotFound - курс не найден
	ErrCourseNotFound = errors.New("course not found")

	// ErrCompanyIDNotEmpty - id компании не может быть пустым
	ErrCompanyIDNotEmpty = errors.New("company id cannot be empty")

	// ErrCompanyReference - id компании должен ссылаться на существующую компанию
	ErrCompanyReference = errors.New("company id must reference existing company")

	// ErrPositionIDNotEmpty - id должности не может быть пустым
	ErrPositionIDNotEmpty = errors.New("position id cannot be empty")

	// ErrPositionReference - id должности должен ссылаться на существующую должность
	ErrPositionReference = errors.New("position id must reference existing position")

	// ErrCreaterNotFound - создатель (пользователь) не должен быть пустым
	ErrCreaterNotEmpty = errors.New("creator (user id) cannot be empty")

	// ErrCreaterNotFound - создатель (пользователь) не зарегистрирован в системе
	ErrCreaterNotFound = errors.New("creator (user id) does not exist in system")

	// ErrCourseIDNotEmpty - id курса не должен быть пустым
	ErrCourseIDNotEmpty = errors.New("course id cannot be empty")

	// ErrCourseIDNotEmpty - id курс должен ссылаться на существующий курс
	ErrCourseReference = errors.New("course id must reference existing course")

	// ErrLessonIDNotEmpty - id урока не должен быть пустым
	ErrLessonIDNotEmpty = errors.New("lesson id cannot be empty")

	// ErrLessonReference - id урока должен ссылаться на существующий урок
	ErrLessonReference = errors.New("lesson id must reference existing lesson")

	// ErrUserNotEmpty - сотрудник (user id) не должен быть пустым
	ErrUserIDNotEmpty = errors.New("user (user id) cannot be empty")

	// ErrUserReference - сотрудник (user id) не зарегистрирован в системе
	ErrUserReference = errors.New("user (user id) does not exist in system")
)

var (
	// -- Частные ошибки объектов --

	// -- Ошибки компании --

	// ErrCompanyNameNotEmpty - название компании не может быть пустым
	ErrCompanyNameNotEmpty = errors.New("company name cannot be empty")

	// ErrIncorrectCompanyName - некорректное имя компании
	ErrIncorrectCompanyName = errors.New("incorrect company name")

	// -- Ошибки должности --

	// ErrPositionNameNotEmpty - название должности не может быть пустым
	ErrPositionNameNotEmpty = errors.New("position name cannot be empty")

	// -- Ошибки пользователя --

	// ErrEmailAlreadyExists - email должен быть уникальный
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrEmailNotEmpty - email не может быть пустым
	ErrEmailNotEmpty = errors.New("email cannot be empty")

	// ErrInvalidEmail - некорректный email
	ErrInvalidEmail = errors.New("invalid email")

	// ErrPasswordNotEmpty - password не может быть пустым
	ErrPasswordNotEmpty = errors.New("password cannot be empty")

	// ErrInvalidPassword - невалидный пароль
	ErrInvalidPassword = errors.New("invalid password: password must have length of 6-30, contain 1 uppercase, 1 lowercase, 1 number, and 1 special character")

	// -- Ошибки курсов --

	// ErrCourseUserNotFound - имя курса не должно быть пустым
	ErrCourseNameNotEmpty = errors.New("course name cannot be empty")

	// -- Ошибки уроков --

	// ErrLessonNameNotEmpty - название урока не может быть пустым
	ErrLessonNameNotEmpty = errors.New("lesson name cannot be empty")

	// -- Ошибки текстов --

	ErrTextContentNotEmpty = errors.New("text (content) cannot be empty")

	// -- Ошибки картинок --

	ErrPictureLinkNotEmpty = errors.New("picture (link) cannot be empty")

	// -- Ошибки назначений и прогресса по учебным материалам --

	// ErrPositionCourseUsed - курс на должность можно назначить только один раз
	ErrPositionCourseUsed = errors.New("course already assigned to position")

	// ErrPositionCourseUsed - курс на сотрудника можно назначить только один раз
	ErrUserCourseUsed = errors.New("course already assigned to user")

	// ErrAssignLessonUsed - прогресс пользователя по уроку должен быть уникальным
	ErrAssignLessonUsed = errors.New("Course progress already has the same lesson progress")
)
