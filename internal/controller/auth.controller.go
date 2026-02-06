package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/apperror"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/response"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/service"
	jwtutil "github.com/NugrahaPancaWibisana/backend-social-media/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
//
//	@Summary		Register new user
//	@Description	Create a new user account
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"User registration data"
//	@Success		201		{object}	dto.ResponseSuccess
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/auth/register [post]
func (ac *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Email field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "email") {
			response.Error(ctx, http.StatusBadRequest, "Email must be a valid email address")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Password field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		if strings.Contains(errStr, "ConfirmPassword") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Confirm password field cannot be empty")
			return
		}

		if strings.Contains(errStr, "ConfirmPassword") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Confirm password must be at least 8 characters")
			return
		}

		if strings.Contains(errStr, "ConfirmPassword") && strings.Contains(errStr, "eqfield") {
			response.Error(ctx, http.StatusBadRequest, "Your password and confirmation password do not match.")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err := ac.authService.Register(ctx, req)

	if err != nil {
		if errors.Is(err, apperror.ErrEmailAlreadyExists) || errors.Is(err, apperror.ErrInvalidEmailFormat) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusCreated, "Registration successful", nil)
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user with email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	dto.ResponseSuccess
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		401		{object}	dto.ResponseError
//	@Router			/auth/login [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Email field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "email") {
			response.Error(ctx, http.StatusBadRequest, "Email must be a valid email address")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Password field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	data, err := ac.authService.Login(ctx, req)
	if err != nil {
		if errors.Is(err, apperror.ErrInvalidEmailFormat) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, apperror.ErrUserNotFound) || errors.Is(err, apperror.ErrInvalidCredential) {
			response.Error(ctx, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	token, err := ac.authService.GenerateJWT(ctx, data)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	ac.authService.WhitelistToken(ctx, data.ID, token)

	response.Success(ctx, http.StatusOK, "Login successful", dto.JWT{Token: token})
}

// Logout godoc
//
//	@Summary		User logout
//	@Description	Logout user and invalidate token
//	@Tags			auth
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	dto.ResponseSuccess
//	@Failure		401	{object}	dto.ResponseError
//	@Failure		500	{object}	dto.ResponseError
//	@Router			/auth/login [delete]
//	@Security		BearerAuth
func (ac *AuthController) Logout(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")
	if len(token) != 2 {
		response.Error(ctx, http.StatusUnauthorized, "Invalid Token")
		return
	}
	if token[0] != "Bearer" {
		response.Error(ctx, http.StatusUnauthorized, "Invalid Token")
		return
	}

	tokenData, _ := ctx.Get("token")
	accessToken, _ := tokenData.(jwtutil.JwtClaims)

	err := ac.authService.Logout(ctx, accessToken.UserID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusOK, "Logout successful", nil)
}
