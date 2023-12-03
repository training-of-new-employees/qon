package jwttoken

import "time"

type JWTGenerator interface {
	GenerateToken(id int, email string, isAdmin bool, orgID int, exp time.Duration) (string, error)
}

type JWTValidator interface {
	ValidateToken(tokenStr string) (*MyClaims, error)
}
