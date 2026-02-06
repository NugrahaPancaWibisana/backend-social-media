package repository

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/apperror"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/model"
	"github.com/jackc/pgx/v5"
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

func (ar *AuthRepository) Login(ctx context.Context, db DBTX, email string) (model.Account, error) {
	sql := `
		SELECT
		    id,
		    email,
			password
		FROM
			accounts
		WHERE 
			email = $1 
			AND deleted_at IS NULL;
	`

	row := db.QueryRow(ctx, sql, email)

	var user model.Account
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		log.Println("ERROR [repostory:auth] failed to login:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Account{}, apperror.ErrUserNotFound
		}
		return model.Account{}, err
	}

	return user, nil
}

func (ar *AuthRepository) UpdateLastLogin(ctx context.Context, db DBTX, id int) error {
	sql := `
		UPDATE 
			accounts
		SET
		    lastlogin_at = NOW()
		WHERE
		    id = $1;
	`

	_, err := db.Exec(ctx, sql, id)
	if err != nil {
		log.Println("ERROR [repostory:auth] failed to update last login:", err)
		return err
	}

	return nil
}
