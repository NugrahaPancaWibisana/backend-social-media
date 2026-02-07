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

func PostRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	postRouter := app.Group("/posts")
	postRouter.Use(middleware.AuthMiddleware())

	postRepository := repository.NewPostRepository()
	postService := service.NewPostService(postRepository, rdb, db)
	postController := controller.NewPostController(postService)

	postRouter.POST("", postController.CreatePost)
	postRouter.GET("/feed", postController.GetFeedPosts)
	postRouter.POST("/like", postController.CreateLike)
	postRouter.POST("/comment", postController.CreateComment)
}
