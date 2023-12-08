package model

import (
	"golang.org/x/crypto/bcrypt"
)

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
