package jwttoken

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

var _ JWTValidator = (*TokenValidator)(nil)

// TokenValidator имплементирует JWTValidator.
type TokenValidator struct {
	secretKey string
}

func NewTokenValidator(secretKey string) *TokenValidator {
	return &TokenValidator{secretKey: secretKey}
}

func (t *TokenValidator) ValidateToken(tokenStr string) (*MyClaims, error) {
	myClaims := MyClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &myClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error failed ValidateToken %w", err)
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
