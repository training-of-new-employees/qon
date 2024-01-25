package model

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/training-of-new-employees/qon/internal/errs"

	"github.com/mcnijman/go-emailaddress"
)

var errSpaceEmpty = errors.New("string only contains spaces")

// validateEmail - проверка емейла.
// ВАЖНО: используется при валидации с методами пакета ozzo-validation.
func validateEmail(email *string) validation.RuleFunc {
	return func(value interface{}) error {
		var err error

		// проверка емейла на корректность
		if err := validation.Validate(email, validation.Length(7, 50), is.Email); err != nil {
			return errs.ErrInvalidEmail
		}

		*email, err = modifyEmail(*email)
		if err != nil {
			return errs.ErrInvalidEmail
		}

		return nil
	}
}

// modifyEmail - преобразование емейла.
func modifyEmail(email string) (string, error) {
	// email должен быть регистронезависимым
	email = strings.ToLower(email)

	emailObj, err := emailaddress.Parse(email)
	if err != nil {
		return "", err
	}
	if emailObj.Domain == "ya.ru" {
		email = emailObj.LocalPart + "@" + "yandex.ru"
	}

	return email, nil
}

// validatePassword - проверка пароля на состав.
// ВАЖНО: используется при валидации с методами пакета ozzo-validation.
func validatePassword(password string) validation.RuleFunc {
	return func(value interface{}) error {
		// Минимум 1 цифра
		numeric := regexp.MustCompile(`\d`).MatchString(password)
		// Минимум 1 буква в нижнем регистре
		lowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
		// Минимум 1 буква в верхнем регистре
		uppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
		// Минимум 1 специальный символ
		special := strings.ContainsAny(password, "!@#$%^&*()_.+")

		only := regexp.MustCompile(`^[0-9a-zA-Z!@#$%^&*()_.+]+$`).MatchString(password)

		if !(numeric && lowercase && uppercase && special && only) {
			return errs.ErrInvalidPassword
		}

		return nil
	}
}

// validateUserName - проверка имени, отчества, фамилии на состав.
// ВАЖНО: используется при валидации с методами пакета ozzo-validation.
func validateUserName(str *string) validation.RuleFunc {
	return func(value interface{}) error {
		// случай когда строка состоит только из пробелов
		*str = strings.TrimSpace(*str)
		if *str == "" {
			return errSpaceEmpty
		}
		for _, c := range *str {
			if !unicode.IsLetter(c) && c != '-' {
				return errors.New("string may only contain unicode characters and a dash")
			}
		}
		return nil
	}
}

// validateNameDescription - проверка имени и описания объектов на состав (компания, должность, курс, урок).
// ВАЖНО: используется при валидации с методами пакета ozzo-validation.
func validateNameDescription(str *string) validation.RuleFunc {
	return func(value interface{}) error {
		// случай когда строка состоит только из пробелов
		*str = strings.TrimSpace(*str)
		if *str == "" {
			return errSpaceEmpty
		}
		for _, c := range *str {
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsPunct(c) &&
				c != '!' && c != '№' && c != ':' && c != '"' && c != '\'' && c != '&' &&
				c != '-' && c != '+' && c != ' ' ||
				c == '*' || c == '#' {
				return errors.New("string may only contain unicode characters and not contain '#' and '*'")
			}
		}

		return nil
	}
}
