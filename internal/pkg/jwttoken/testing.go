package jwttoken

import (
	"time"
)

// TestAuthorizateUser - получение тестового токена для авторизации пользователя.
// Cценарии для админа и сотрудника конфигурируются с помощью аргумента isAdmin.
func TestAuthorizateUser(userID int, companyID int, isAdmin bool) (string, error) {
	g := NewTokenGenerator("secret")

	accessToken, err := g.GenerateToken(userID, isAdmin, companyID, "", time.Minute*10)
	if err != nil {
		return "", err
	}

	return "Bearer " + accessToken, nil
}
