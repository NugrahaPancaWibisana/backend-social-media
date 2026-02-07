package controller

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/response"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/service"
	jwtutil "github.com/NugrahaPancaWibisana/backend-social-media/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type PostController struct {
	postService *service.PostService
}

func NewPostController(postService *service.PostService) *PostController {
	return &PostController{postService: postService}
}

// CreatePost godoc
//
//	@Summary		Create new post
//	@Description	Create a new post
//	@Tags			Posts
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			content		formData	file	false	"content image"
//	@Param			caption		formData	string	false	"caption (min 3 chars)"
//	@Success		201		{object}	dto.ResponseSuccess
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/posts [post]
//	@Security		BearerAuth
func (pc *PostController) CreatePost(ctx *gin.Context) {
	var req dto.PostRequest

	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "no multipart boundary param in Content-Type") {
			response.Error(ctx, http.StatusBadRequest, "Invalid multipart form data")
			return
		}

		if strings.Contains(errStr, "Caption") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Caption must be at least 3 characters")
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
	if req.Content != nil {
		ext := strings.ToLower(path.Ext(req.Content.Filename))
		re := regexp.MustCompile(`^\.(jpg|png)$`)
		if !re.MatchString(ext) {
			response.Error(ctx, http.StatusBadRequest, "File must be jpg or png")
			return
		}

		filename := fmt.Sprintf(
			"%d_post_%d%s",
			time.Now().UnixNano(),
			accessToken.UserID,
			ext,
		)

		if err := ctx.SaveUploadedFile(
			req.Content,
			filepath.Join("public", "post", filename),
		); err != nil {
			log.Println(err.Error())
			response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		imagePath = fmt.Sprintf("/post/%s", filename)
	}

	if err := pc.postService.CreatePost(ctx, req, accessToken.UserID, imagePath, token[1]); err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusOK, "Create post successful", nil)
}
