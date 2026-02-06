package router

import (
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/controller"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/middleware"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/repository"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func UserRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	userRouter := app.Group("/users")
	userRouter.Use(middleware.AuthMiddleware())

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, rdb, db)
	userController := controller.NewUserController(userService)

	userRouter.GET("/profile", userController.GetProfile)
}
