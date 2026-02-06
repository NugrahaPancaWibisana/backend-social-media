package dto

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID int    `json:"id"`
	jwt.RegisteredClaims
}

type JWT struct {
	Token string `json:"token"`
}

type Account struct {
	ID          int        `json:"id" example:"1"`
	Email       string     `json:"email" example:"user@example.com"`
}
