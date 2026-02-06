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

type UserService struct {
	userRepository *repository.UserRepository
	redis          *redis.Client
	db             *pgxpool.Pool
}

func NewUserService(userRepository *repository.UserRepository, rdb *redis.Client, db *pgxpool.Pool) *UserService {
	return &UserService{userRepository: userRepository, redis: rdb, db: db}
}

func (us *UserService) GetProfile(ctx context.Context, id int, token string) (dto.User, error) {
	err := cache.CheckToken(ctx, us.redis, id, token)
	if err != nil {
		log.Println("ERROR [service:user] failed to get profile:", err)
		return dto.User{}, err
	}

	data, err := us.userRepository.GetProfile(ctx, us.db, id)
	if err != nil {
		log.Println("ERROR [service:user] failed to get profile:", err)
		return dto.User{}, err
	}

	res := dto.User{
		ID:     data.ID,
		Email:  data.Email,
		Name:   data.Name.String,
		Avatar: data.Avatar.String,
		Bio:    data.Bio.String,
	}

	return res, nil
}

func (us *UserService) UpdateProfile(ctx context.Context, req dto.UpdateProfileRequest, path string, id int, token string) (string, error) {
	err := cache.CheckToken(ctx, us.redis, id, token)
	if err != nil {
		log.Println("ERROR [service:user] failed to update profile:", err)
		return "", err
	}

	tx, err := us.db.Begin(ctx)
	if err != nil {
		log.Println("ERROR [service:user] failed to begin:", err)
		return "", err
	}
	defer tx.Rollback(ctx)

	oldPath, err := us.userRepository.GetAvatar(ctx, tx, id)
	if err != nil {
		log.Println("ERROR [service:user] failed to get avatar:", err)
		return "", err
	}

	if err := us.userRepository.UpdateProfile(ctx, tx, req, path, id); err != nil {
		log.Println("ERROR [service:user] failed to get avatar:", err)
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("ERROR [service:user] failed to commit:", err)
		return "", err
	}

	return oldPath.String, nil
}
