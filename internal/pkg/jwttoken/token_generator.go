package jwttoken

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var _ JWTGenerator = (*TokenGenerator)(nil)

// TokenGenerator имплементирует JWTGenerator.
type TokenGenerator struct {
	secretKey string
	//accessToken  time.Duration
	//refreshToken time.Duration
}

type MyClaims struct {
	jwt.RegisteredClaims
	UserID  int    `json:"id"`
	Email   string `json:"email"`
	ReqID   string `json:"req_id"`
	IsAdmin bool   `json:"is_admin"`
	OrgID   int    `json:"org_id"`
}

func NewTokenGenerator(secretKey string) *TokenGenerator {
	return &TokenGenerator{
		secretKey: secretKey,
	}
}

func (g *TokenGenerator) GenerateToken(id int, email string, isAdmin bool, orgID int, exp time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := MyClaims{
		UserID:  id,
		Email:   email,
		IsAdmin: isAdmin,
		OrgID:   orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}

	token.Claims = claims

	tokenStr, err := token.SignedString([]byte(g.secretKey))
	if err != nil {
		return "", fmt.Errorf("err failed GenerateToken %v", err)
	}

	return tokenStr, nil
}
