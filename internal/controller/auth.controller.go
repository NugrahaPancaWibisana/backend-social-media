package controller

import "github.com/NugrahaPancaWibisana/backend-social-media/internal/service"

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}