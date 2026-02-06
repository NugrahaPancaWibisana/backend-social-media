package repository

import (
	"context"
	"log"
	"strings"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/apperror"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (ar *AuthRepository) CreateAccount(ctx context.Context, db DBTX, req dto.RegisterRequest) error {
	sql := `
		INSERT INTO
		    accounts (email, password)
		VALUES
		    ($1, $2)
	`

	if _, err := db.Exec(ctx, sql, req.Email, req.Password); err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "duplicate") {
			return apperror.ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}
