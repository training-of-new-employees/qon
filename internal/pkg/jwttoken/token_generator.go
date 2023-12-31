package jwttoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var _ JWTGenerator = (*TokenGenerator)(nil)

// TokenGenerator имплементирует JWTGenerator.
type TokenGenerator struct {
	secretKey string
}

type MyClaims struct {
	jwt.RegisteredClaims
	UserID  int  `json:"id"`
	IsAdmin bool `json:"is_admin"`
	OrgID   int  `json:"org_id"`
}

func NewTokenGenerator(secretKey string) *TokenGenerator {
	return &TokenGenerator{
		secretKey: secretKey,
	}
}

func (g *TokenGenerator) GenerateToken(id int, isAdmin bool, orgID int, exp time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := MyClaims{
		UserID:  id,
		IsAdmin: isAdmin,
		OrgID:   orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}

	token.Claims = claims

	tokenStr, err := token.SignedString([]byte(g.secretKey))
	if err != nil {
		return "", fmt.Errorf("error failed GenerateToken %v", err)
	}

	return tokenStr, nil
}
