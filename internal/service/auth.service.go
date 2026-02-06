package service

import (
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	authRepository *repository.AuthRepository
	redis          *redis.Client
	db             *pgxpool.Pool
}

func NewAuthService(authRepository *repository.AuthRepository, rdb *redis.Client, db *pgxpool.Pool) *AuthService {
	return &AuthService{authRepository: authRepository, redis: rdb, db: db}
}