package dto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID int    `json:"id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWT struct {
	Token string `json:"token"`
}

type Account struct {
	ID          int        `json:"id" example:"1"`
	Email       string     `json:"email" example:"user@example.com"`
	LastLoginAt *time.Time `json:"lastlogin_at" example:"2026-02-06T14:30:45+01:00"`
}
