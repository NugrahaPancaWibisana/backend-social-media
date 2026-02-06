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

func (ar *AuthRepository) Login(ctx context.Context, db DBTX, email string) (model.User, error) {
	query := `
		SELECT
		    id,
		    email,
			password,
			lastlogin_at
		FROM
			accounts
		WHERE 
			email = $1 
			AND deleted_at IS NULL;
	`

	row := db.QueryRow(ctx, query, email)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.LastLoginAt,
	)

	if err != nil {
		log.Println("ERROR [repostory:auth] failed to login:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, apperror.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (ar *AuthRepository) UpdateLastLogin(ctx context.Context, db DBTX, id int) error {
	query := `
		UPDATE 
			users
		SET
		    lastlogin_at = NOW()
		WHERE
		    id = $1;
	`

	_, err := db.Exec(ctx, query, id)
	if err != nil {
		log.Println("ERROR [repostory:auth] failed to update last login:", err)
		return err
	}

	return nil
}
