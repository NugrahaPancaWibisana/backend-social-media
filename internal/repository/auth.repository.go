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

func (ar *AuthRepository) CreateAccount(ctx context.Context, db DBTX, req dto.RegisterRequest) (int, error) {
	sql := `
		INSERT INTO
		    accounts (email, password)
		VALUES
		    ($1, $2)
		RETURNING id;
	`

	var id int
	err := db.QueryRow(ctx, sql, req.Email, req.Password).Scan(&id)
	if err != nil {
		log.Println("ERROR [repostory:auth] failed to create account:", err)
		if strings.Contains(err.Error(), "duplicate") {
			return 0, apperror.ErrEmailAlreadyExists
		}
		return 0, err
	}

	return id, nil
}

func (ar *AuthRepository) CreateUser(ctx context.Context, db DBTX, id int) error {
	sql := `
		INSERT INTO
		    users (account_id)
		VALUES
		    ($1);
	`

	if _, err := db.Exec(ctx, sql, id); err != nil {
		log.Println("ERROR [repostory:auth] failed to create user profile:", err)
		return err
	}

	return nil
}
