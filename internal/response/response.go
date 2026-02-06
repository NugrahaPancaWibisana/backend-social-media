package response

import (
	"net/http"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, statusCode int, message string, data any) {
	res := dto.Response{
		Status:  "success",
		Message: message,
	}

	if data != nil {
		ctx.JSON(statusCode, dto.ResponseSuccess{
			Response: res,
			Data:     data,
		})
		return
	}

	ctx.JSON(statusCode, res)
}

func SuccessWithMeta(ctx *gin.Context, statusCode int, message string, data any, meta any) {
	ctx.JSON(statusCode, dto.ResponseSuccessWithMeta{
		Response: dto.Response{
			Status:  "success",
			Message: message,
		},
		Data: data,
		Meta: meta,
	})
}

func Error(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, dto.ResponseError{
		Response: dto.Response{
			Status:  "error",
			Message: message,
		},
		Error: http.StatusText(statusCode),
	})
}

func Abort(ctx *gin.Context, statusCode int, message string) {
	ctx.AbortWithStatusJSON(statusCode, dto.ResponseError{
		Response: dto.Response{
			Status:  "success",
			Message: message,
		},
		Error: http.StatusText(statusCode),
	})
}
