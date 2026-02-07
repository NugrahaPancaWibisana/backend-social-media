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

func (ps *PostService) GetFeedPosts(ctx context.Context, userID int, token string) ([]dto.FeedPost, error) {

	if err := cache.CheckToken(ctx, ps.redis, userID, token); err != nil {
		log.Println("ERROR [service:feed] token invalid:", err)
		return nil, err
	}

	data, err := ps.postRepository.GetFeedPosts(ctx, ps.db, userID)
	if err != nil {
		log.Println("ERROR [service:feed] failed to get feed:", err)
		return nil, err
	}

	res := make([]dto.FeedPost, 0, len(data))
	for _, p := range data {
		res = append(res, dto.FeedPost{
			ID:        p.ID,
			Content:   p.Content,
			Caption:   p.Caption,
			Likes:     p.LikesCount,
			Comments:  p.CommentsCount,
			CreatedAt: p.CreatedAt,
			Author: dto.Author{
				ID:     p.UserID,
				Name:   p.UserName,
				Avatar: p.UserAvatar,
			},
		})
	}

	return res, nil
}
