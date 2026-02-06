package repository

import (
	"context"
	"errors"
	"log"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/apperror"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/model"
	"github.com/jackc/pgx/v5"
)

type UserRepo interface {
	GetProfile(ctx context.Context, db DBTX, id int) (model.User, error)
}

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) GetProfile(ctx context.Context, db DBTX, userID int) (model.User, error) {
	sql := `
		SELECT
			a.id,
			a.email,
			u.name,
			u.avatar,
			u.bio
		FROM
			accounts a
			JOIN users u ON u.account_id = a.id
		WHERE id = $1;
	`

	row := db.QueryRow(ctx, sql, userID)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Avatar,
		&user.Bio,
	)

	if err != nil {
		log.Println("ERROR [repostory:auth] failed to get profile:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, apperror.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}
