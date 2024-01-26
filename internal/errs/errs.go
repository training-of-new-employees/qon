// Package errs - пакет содержит ошибки приложения.
package errs

import (
	"errors"
	"fmt"
)

const specialSymbols = "(),?!№:\"\\&-_%';@"

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
	// ErrUserActivated - пользователь уже активирован в системе
	ErrUserActivated = errors.New("user is activated")
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
	// ErrNoAccess - нет доступа
	ErrNoAccess = errors.New("no access")
	// ErrIncorrectEmailOrPassword - неправильный емейл или пароль
	ErrIncorrectEmailOrPassword = errors.New("incorrect email or password")
	// ErrEmailOrPasswordEmpty - пустой емейл или пароль
	ErrEmailOrPasswordEmpty = errors.New("email address and password must be filled in")
	// ErrVerifyCodeNotEmpty - пустой код верификации
	ErrVerifyCodeNotEmpty = errors.New("code cannot be empty")
	// ErrIncorrectVerifyCode - невалидный код верификации
	ErrIncorrectVerifyCode = errors.New("invalid verify code")
	// ErrInvalidInviteCode - невалидный пригласительный код
	ErrInvalidInviteCode = errors.New("something wrong with invite code: code is invalid or something else")
	// ErrInvalidRoute - не найденный маршрут
	ErrInvalidRoute = errors.New("route not found")
)

var (
	// -- Общие ошибки объектов --

	// ErrCompanyNotFound - компания не найдена
	ErrCompanyNotFound = errors.New("company not found")
	// ErrPositionNotFound - должность не найдена
	ErrPositionNotFound = errors.New("position not found")
	// ErrCompanyNoPosition - компания и должность не имеют связи
	ErrCompanyNoPosition = errors.New("position and company are not related")
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
	ErrInvalidCompanyName = fmt.Errorf("invalid company name: company name must have length of 1-256 and can contain characters of any alphabets, digits, spaces, special symbols %s", specialSymbols)

	// -- Ошибки должности --

	// ErrPositionNameNotEmpty - название должности не может быть пустым
	ErrPositionNameNotEmpty = errors.New("position name cannot be empty")

	ErrInvalidPositionName = fmt.Errorf("invalid position name: position name must have length of 2-256 and can contain characters of any alphabets, digits, spaces, special symbols %s", specialSymbols)

	// -- Ошибки пользователя --

	// ErrEmailAlreadyExists - email должен быть уникальный
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrEmailNotEmpty - email не может быть пустым
	ErrEmailNotEmpty = errors.New("email cannot be empty")

	// ErrInvalidEmail - некорректный email
	ErrInvalidEmail = errors.New("invalid email or incorrect length (email length should be 7-50 symbols)")

	// ErrPasswordNotEmpty - password не может быть пустым
	ErrPasswordNotEmpty = errors.New("password cannot be empty")

	// ErrInviteNotEmpty - code invite не может быть пустым
	ErrInviteNotEmpty = errors.New("invite cannot be empty")

	// ErrInvalidPassword - невалидный пароль
	ErrInvalidPassword = errors.New("invalid password: password must have length of 6-30, contain 1 uppercase, 1 lowercase, 1 number, and 1 special character")

	// ErrUserNameNotEmpty - имя сотрудника не должно быть пустым
	ErrUserNameNotEmpty = errors.New("user name cannot be empty")

	// ErrIncorrectUserName - некорректное имя пользователя
	ErrInvalidUserName = errors.New("invalid user name: name must have length of 2-128 and can contain characters of any alphabets, dash")

	// ErrUserSurnameNotEmpty - фамилия сотрудника не должна быть пустой
	ErrUserSurnameNotEmpty = errors.New("user surname cannot be empty")

	// ErrInvalidUserSurname - некорректная фамилия пользователя
	ErrInvalidUserSurname = errors.New("invalid user surname: surname must have length of 2-128 and can contain characters of any alphabets, dash")

	// ErrUserPatronymicNotEmpty - отчество сотрудника не должна быть пустой
	ErrUserPatronymicNotEmpty = errors.New("user patronymic cannot be empty")

	// ErrInvalidUserPatronymic - некорректное отчество пользователя
	ErrInvalidUserPatronymic = errors.New("invalid user patronymic: patronymic must have length of 2-128 and can contain characters of any alphabets, dash")

	// -- Ошибки курсов --

	// ErrCourseUserNotFound - имя курса не должно быть пустым
	ErrCourseNameIsEmpty = errors.New("course name cannot be empty")

	// ErrInvalidCourseName - имя курса не соответствует требованиям
	ErrInvalidCourseName = fmt.Errorf("course name is invalid: course name must have length of 5-256 and can contain characters of any alphabets, digits, spaces, special symbols %s", specialSymbols)

	// ErrInvalidCourseDescription - описание курса не соответствует требованиям
	ErrInvalidCourseDescription = fmt.Errorf("course description is invalid: course name is invalid: course description must have length of 10-512 and can contain characters of any alphabets, digits, spaces, special symbols %s", specialSymbols)

	// -- Ошибки уроков --

	// ErrLessonNameNotEmpty - название урока не может быть пустым
	ErrLessonNameNotEmpty = errors.New("lesson name cannot be empty")

	// ErrInvalidLessonName - имя урока не соответствует требованиям
	ErrInvalidLessonName = fmt.Errorf("invalid lesson name: lesson name must have length of 5-256 and can contain characters of any alphabets, digits, spaces, special symbols %s", specialSymbols)

	// -- Ошибки текстов --

	ErrTextContentNotEmpty = errors.New("text (content) cannot be empty")

	ErrInvalidTextContent = fmt.Errorf("invalid lesson content: lesson content must have length of 20-65000 and can contain characters of any alphabets, digits, spaces, special symbols %s", specialSymbols)

	// -- Ошибки картинок --

	ErrURLPictureNotEmpty = errors.New("picture's url cannot be empty")

	ErrURLPictureLength = errors.New("picture's url must have length of 5-1024")

	// -- Ошибки назначений и прогресса по учебным материалам --

	// ErrPositionCourseUsed - курс на должность можно назначить только один раз
	ErrPositionCourseUsed = errors.New("course already assigned to position")

	// ErrPositionCourseUsed - курс на сотрудника можно назначить только один раз
	ErrUserCourseUsed = errors.New("course already assigned to user")

	// ErrAssignLessonUsed - прогресс пользователя по уроку должен быть уникальным
	ErrAssignLessonUsed = errors.New("Course progress already has the same lesson progress")

	// ErrInvalidCourseStatus - невалидный статус сотрудника по курсу
	ErrInvalidCourseStatus = errors.New("user course status can be 'not-started', 'in-process', 'done'")

	// ErrInvalidLessonStatus - невалидный статус сотрудника по уроку
	ErrInvalidLessonStatus = errors.New("user lesson status can be 'not-started', 'in-process', 'done'")
)
