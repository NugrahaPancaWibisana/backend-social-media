package controller

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/apperror"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/response"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/service"
	jwtutil "github.com/NugrahaPancaWibisana/backend-social-media/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
//	@Failure		500	{object}	dto.ResponseError
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

// UpdateProfile godoc
//
//	@Summary		Update user profile
//	@Description	Update authenticated user's profile information
//	@Tags			Users
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			avatar		formData	file	false	"Profile avatar"
//	@Param			name		formData	string	false	"Name (min 3 chars)"
//	@Param			bio			formData	string	false	"bio (min 3 chars)"
//	@Success		200			{object}	dto.ResponseSuccess
//	@Failure		400			{object}	dto.ResponseError
//	@Failure		401			{object}	dto.ResponseError
//	@Failure		500	{object}	dto.ResponseError
//	@Router			/users/profile [patch]
//	@Security		BearerAuth
func (uc *UserController) UpdateProfile(ctx *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "no multipart boundary param in Content-Type") {
			response.Error(ctx, http.StatusBadRequest, "No fields to update")
			return
		}

		if strings.Contains(errStr, "Avatar") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Avatar must be at least 3 characters")
			return
		}

		if strings.Contains(errStr, "Name") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Name must be at least 3 characters")
			return
		}

		if strings.Contains(errStr, "Bio") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Bio must be at least 3 characters")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

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

	var imagePath string

	if req.Avatar != nil {
		const maxFileSize = 2 * 1024 * 1024
		if req.Avatar.Size > maxFileSize {
			response.Error(ctx, http.StatusBadRequest, "File size must not exceed 2 MB")
			return
		}

		ext := strings.ToLower(path.Ext(req.Avatar.Filename))
		re := regexp.MustCompile(`^\.(jpg|png)$`)
		if !re.MatchString(ext) {
			response.Error(ctx, http.StatusBadRequest, "File must be jpg or png")
			return
		}

		file, err := req.Avatar.Open()
		if err != nil {
			log.Println("failed to open uploaded file:", err.Error())
			response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		defer file.Close()

		img, _, err := image.DecodeConfig(file)
		if err != nil {
			log.Println("failed to decode image:", err.Error())
			response.Error(ctx, http.StatusBadRequest, "Invalid image file")
			return
		}

		const maxWidth = 800
		const maxHeight = 600
		if img.Width > maxWidth || img.Height > maxHeight {
			response.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Image dimensions must not exceed %dx%d pixels", maxWidth, maxHeight))
			return
		}

		filename := fmt.Sprintf(
			"%d_profile_%d%s",
			time.Now().UnixNano(),
			accessToken.UserID,
			ext,
		)

		if err := ctx.SaveUploadedFile(
			req.Avatar,
			filepath.Join("public", "profile", filename),
		); err != nil {
			log.Println(err.Error())
			response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		imagePath = fmt.Sprintf("/profile/%s", filename)
	}

	oldPath, err := uc.userService.UpdateProfile(ctx, req, imagePath, accessToken.UserID, token[1])
	if err != nil {
		if errors.Is(err, apperror.ErrNoFieldsToUpdate) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if oldPath != "" && imagePath != "" {
		oldFullPath := filepath.Join("public", oldPath)
		if err := os.Remove(oldFullPath); err != nil {
			log.Println("failed to delete old photo:", err.Error())
		}
	}

	response.Success(ctx, http.StatusOK, "Profile updated successfully", nil)
}

// GetUsers godoc
//
//	@Summary		Get users
//	@Description	Get list of users
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	dto.Response
//	@Failure		401	{object}	dto.ResponseError
//	@Failure		500	{object}	dto.ResponseError
//	@Router			/users [get]
//	@Security		BearerAuth
func (uc *UserController) GetUsers(ctx *gin.Context) {
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

	data, err := uc.userService.GetUsers(ctx, accessToken.UserID, token[1])
	if err != nil {
		response.Error(
			ctx,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
		return
	}

	response.Success(ctx, http.StatusOK, "Users retrieved successfully", data)
}

// FollowUser godoc
//
//	@Summary		Follow user
//	@Description	Follow another user
//	@Tags			Users
//	@Produce		json
//	@Param			id	path	int	true	"User ID to follow"
//	@Success		200	{object}	dto.ResponseSuccess
//	@Failure		400	{object}	dto.ResponseError
//	@Failure		401	{object}	dto.ResponseError
//	@Failure		500	{object}	dto.ResponseError
//	@Router			/users/{id}/follow [post]
//	@Security		BearerAuth
func (uc *UserController) FollowUser(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")
	if len(token) != 2 {
		response.Error(ctx, http.StatusUnauthorized, "Invalid Token")
		return
	}
	if token[0] != "Bearer" {
		response.Error(ctx, http.StatusUnauthorized, "Invalid Token")
		return
	}

	followedID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user id")
		return
	}

	tokenData, _ := ctx.Get("token")
	accessToken, _ := tokenData.(jwtutil.JwtClaims)

	err = uc.userService.FollowUser(
		ctx,
		accessToken.UserID,
		followedID,
		token[1],
	)
	if err != nil {
		if errors.Is(err, apperror.ErrCannotFollowYourself) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(
			ctx,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
		return
	}

	response.Success(ctx, http.StatusOK, "User followed successfully", nil)
}
