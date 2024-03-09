package model

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/mcnijman/go-emailaddress"

	"github.com/training-of-new-employees/qon/internal/errs"
)

var errSpaceEmpty = errors.New("string only contains spaces")

// validateEmail - проверка емейла.
// ВАЖНО: используется при валидации с методами пакета ozzo-validation.
func validateEmail(email *string) validation.RuleFunc {
	return func(value interface{}) error {
		var err error

		// проверка емейла на корректность
		if err := validation.Validate(email, is.Email); err != nil {
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

// validateURLPicture - проверка url-картинки.
// ВАЖНО: используется при валидации с методами пакета ozzo-validation.
func validateURLPicture(url *string) validation.RuleFunc {
	return func(value interface{}) error {
		// проверка url на корректность
		if err := validation.Validate(url, is.URL); err != nil {
			return errs.ErrInvalidURLPicture
		}

		picture := regexp.MustCompile(`.(png|jpg|jpeg)$`).MatchString(*url)
		if !picture {
			return errs.ErrInvalidURLPicture
		}

		return nil
	}
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
		// Только символы 0-9a-zA-Z!@#$%^&*()_.+
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
			if !unicode.IsLetter(c) && c != '-' && c != ' ' && c != '\'' {
				return errors.New("string may only contain unicode characters, dash, space and apostrophe")
			}
		}
		return nil
	}
}

// validateObjName - проверка на состав название объекта (компания, должность, курс, урок).
// ВАЖНО: используется при валидации с методами пакета ozzo-validation.
func validateObjName(str *string) validation.RuleFunc {
	return func(value interface{}) error {
		// случай когда строка состоит только из пробелов
		*str = strings.TrimSpace(*str)
		if *str == "" {
			return errSpaceEmpty
		}
		for _, c := range *str {
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsPunct(c) &&
				c != '!' && c != '№' && c != ':' && c != '"' && c != '\'' && c != '&' &&
				c != '-' && c != '_' && c != '+' && c != ' ' ||
				c == '*' || c == '#' {
				return errors.New("string may only contain unicode characters and not contain '#' and '*'")
			}
		}

		return nil
	}
}

// validateObjDescription - проверка на состав (компания, должность, курс, урок).
// ВАЖНО: используется при валидации с методами пакета ozzo-validation.
func validateObjDescription(str *string) validation.RuleFunc {
	return func(value interface{}) error {
		// случай когда строка состоит только из пробелов
		*str = strings.TrimSpace(*str)
		if *str == "" {
			return errSpaceEmpty
		}
		for _, c := range *str {
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsPunct(c) &&
				c != '!' && c != '№' && c != ':' && c != '"' && c != '\'' && c != '&' &&
				c != '-' && c != '_' && c != '+' && c != ' ' &&
				c != '*' && c != '#' {
				return errors.New("string may only contain unicode characters and special symbols")
			}
		}

		return nil
	}
}
