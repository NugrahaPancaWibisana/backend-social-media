package service

import (
	"context"
	"log"
	"regexp"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/apperror"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
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

func (as *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, req.Email)
	if !matched {
		return apperror.ErrInvalidEmailFormat
	}

	tx, err := as.db.Begin(ctx)
	if err != nil {
		log.Println("ERROR [service:auth] failed to begin:", err)
		return err
	}
	defer tx.Rollback(ctx)

	id, err := as.authRepository.CreateAccount(ctx, tx, req)
	if err != nil {
		log.Println("ERROR [service:auth] failed to create account:", err)
		return err
	}

	if err := as.authRepository.CreateUser(ctx, tx, id); err != nil {
		log.Println("ERROR [service:auth] failed to create user profile:", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("ERROR [service:auth] failed to commit:", err)
		return err
	}

	return nil
}
