package model

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/training-of-new-employees/qon/internal/errs"
)

var errSpaceEmpty = errors.New("string only contains spaces")

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
func validateUserName(str string) validation.RuleFunc {
	return func(value interface{}) error {
		for _, c := range str {
			if !unicode.IsLetter(c) && c != '-' {
				return errors.New("string may only contain unicode characters and a dash")
			}
		}
		return nil
	}
}

// validateNameDescription - проверка имени и описания объектов на состав (компания, должность, курс, урок).
// ВАЖНО: используется при валидации с методами пакета ozzo-validation.
func validateNameDescription(str string) validation.RuleFunc {
	return func(value interface{}) error {
		// случай когда строка состоит только из пробелов
		trimmed := strings.Trim(str, " ")
		if trimmed == "" {
			return errSpaceEmpty
		}
		for _, c := range str {
			// only (),?!№:"\&-_%';@
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
