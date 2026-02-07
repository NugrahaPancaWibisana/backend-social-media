package service

import (
	"context"
	"log"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/cache"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type PostService struct {
	postRepository *repository.PostRepository
	redis          *redis.Client
	db             *pgxpool.Pool
}

func NewPostService(postRepository *repository.PostRepository, rdb *redis.Client, db *pgxpool.Pool) *PostService {
	return &PostService{postRepository: postRepository, redis: rdb, db: db}
}

func (ps *PostService) CreatePost(ctx context.Context, req dto.PostRequest, id int, path, token string) error {
	err := cache.CheckToken(ctx, ps.redis, id, token)
	if err != nil {
		log.Println("ERROR [service:post] failed to create post:", err)
		return err
	}

	if err := ps.postRepository.CreatePost(ctx, ps.db, req, path, id); err != nil {
		log.Println("ERROR [service:post] failed to create post:", err)
		return err
	}

	return nil
}
