package service

import (
	"context"

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
		return dto.User{}, err
	}

	data, err := us.userRepository.GetProfile(ctx, us.db, id)
	if err != nil {
		return dto.User{}, err
	}

	var (
		Name   string
		Avatar string
		Bio    string
	)

	switch {
	case !data.Name.Valid:
		Name = ""
	case !data.Avatar.Valid:
		Avatar = ""
	case !data.Bio.Valid:
		Bio = ""
	}

	res := dto.User{
		ID:     data.ID,
		Email:  data.Email,
		Name:   Name,
		Avatar: Avatar,
		Bio:    Bio,
	}

	return res, nil
}
