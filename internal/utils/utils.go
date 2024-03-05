package utils

import (
	"strconv"

	"github.com/training-of-new-employees/qon/internal/errs"
	"golang.org/x/crypto/bcrypt"
)

// maxIntPostgres - максимальное значение типа int в postgres.
// Подробнее о числовых типах в postgres: https://www.postgresql.org/docs/current/datatype-numeric.html#DATATYPE-NUMERIC
const maxIntPostgres = 2147483647

// ConvertID используется для преобразования строкового представления id в числовое.
func ConvertID(strID string) (int, error) {
	userID, err := strconv.Atoi(strID)
	if err != nil {
		return 0, errs.ErrBadRequest
	}

	if userID <= 0 || userID > maxIntPostgres {
		return 0, errs.ErrBadRequest
	}

	return userID, nil
}

// EncryptPassword используется для хеширования пароля.
func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ComparePassword сравнивает хешированный пароль с его возможным эквивалентом в виде текста.
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
