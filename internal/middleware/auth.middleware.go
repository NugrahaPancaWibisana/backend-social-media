package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/response"
	jwtutil "github.com/NugrahaPancaWibisana/backend-social-media/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := strings.Split(ctx.GetHeader("Authorization"), " ")
		if len(token) != 2 {
			response.Abort(ctx, http.StatusUnauthorized, "Invalid Token")
			return
		}
		if token[0] != "Bearer" {
			response.Abort(ctx, http.StatusUnauthorized, "Invalid Token")
			return
		}

		var jc jwtutil.JwtClaims
		_, err := jc.VerifyToken(token[1])
		if err != nil {
			log.Println(err.Error())
			if errors.Is(err, jwt.ErrTokenExpired) {
				response.Abort(ctx, http.StatusUnauthorized, "Expired Token, Please Login Again")
				return
			}
			if errors.Is(err, jwt.ErrTokenInvalidIssuer) {
				response.Abort(ctx, http.StatusUnauthorized, "Invalid Token, Please Login Again")
				return
			}

			response.Abort(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		ctx.Set("token", jc)
		ctx.Next()
	}
}
