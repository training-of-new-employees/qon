package jwttoken

import "time"

type JWTGenerator interface {
	GenerateToken(id int, isAdmin bool, orgID int, hashedRefreshToken string, exp time.Duration) (string, error)
}

type JWTValidator interface {
	ValidateToken(tokenStr string) (*MyClaims, error)
}
