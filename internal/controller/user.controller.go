package controller

import (
	"net/http"
	"strings"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/response"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/service"
	jwtutil "github.com/NugrahaPancaWibisana/backend-social-media/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetProfile godoc
//
//	@Summary		Get user profile
//	@Description	Get authenticated user's profile information
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	dto.Response
//	@Failure		401	{object}	dto.ResponseError
//	@Router			/users/profile [get]
//	@Security		BearerAuth
func (uc *UserController) GetProfile(ctx *gin.Context) {
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
	data, err := uc.userService.GetProfile(ctx, accessToken.UserID, token[1])
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusOK, "Profile retrieved successfully", data)
}
