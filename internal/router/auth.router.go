package router

import (
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/controller"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/repository"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func AuthRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := app.Group("/auth")

	authRepository := repository.NewAuthRepository()
	authService := service.NewAuthService(authRepository, rdb, db)
	authController := controller.NewAuthController(authService)

	authRouter.POST("/register", authController.Register)
	authRouter.POST("/login", authController.Register)
}
