package service

import (
	"context"
	"log"
	"regexp"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/apperror"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/cache"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/repository"
	hashutil "github.com/NugrahaPancaWibisana/backend-social-media/pkg/hash"
	jwtutil "github.com/NugrahaPancaWibisana/backend-social-media/pkg/jwt"
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

	hasher := hashutil.Default()
	hashedPassword, err := hasher.Hash(req.Password)
	if err != nil {
		log.Println("ERROR [service:auth] failed to hash:", err)
		return err
	}

	req.Password = hashedPassword

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

func (as *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.Account, error) {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, req.Email)
	if !matched {
		return dto.Account{}, apperror.ErrInvalidEmailFormat
	}

	tx, err := as.db.Begin(ctx)
	if err != nil {
		log.Println("ERROR [service:auth] failed to begin:", err)
		return dto.Account{}, err
	}
	defer tx.Rollback(ctx)

	data, err := as.authRepository.Login(ctx, tx, req.Email)
	if err != nil {
		log.Println("ERROR [service:auth] failed to login:", err)
		return dto.Account{}, err
	}

	hasher := hashutil.Default()
	isValid, err := hasher.Verify(req.Password, data.Password)
	if err != nil {
		log.Println("ERROR [service:auth] failed to verify password:", err)
		return dto.Account{}, err
	}
	if !isValid {
		log.Println("ERROR [service:auth] failed, invalid credential:", err)
		return dto.Account{}, apperror.ErrInvalidCredential
	}

	err = as.authRepository.UpdateLastLogin(ctx, tx, data.ID)
	if err != nil {
		log.Println("ERROR [service:auth] failed to update last login:", err)
		return dto.Account{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("ERROR [service:auth] failed to commit:", err)
		return dto.Account{}, err
	}

	res := dto.Account{
		ID:          data.ID,
		Email:       data.Email,
	}

	return res, nil
}

func (as *AuthService) GenerateJWT(ctx context.Context, user dto.Account) (string, error) {
	claims := jwtutil.NewJWTClaims(user.ID)
	return claims.GenToken()
}

func (as *AuthService) WhitelistToken(ctx context.Context, id int, token string) {
	cache.SetToken(ctx, as.redis, id, token)
}

func (as *AuthService) Logout(ctx context.Context, userID int) error {
	return cache.DeleteToken(ctx, as.redis, userID)
}
