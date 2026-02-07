package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/apperror"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
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
		WHERE account_id = $1;
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
		log.Println("ERROR [repostory:user] failed to get profile:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, apperror.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) GetAvatar(ctx context.Context, db DBTX, id int) (sql.NullString, error) {
	query := "SELECT avatar FROM users WHERE account_id = $1;"

	row := db.QueryRow(ctx, query, id)

	var avatar sql.NullString
	err := row.Scan(&avatar)

	if err != nil {
		log.Println("ERROR [repostory:user] failed to get avatar:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return sql.NullString{}, nil
		}
		return sql.NullString{}, err
	}

	return avatar, nil
}

func (ur *UserRepository) UpdateProfile(ctx context.Context, db DBTX, req dto.UpdateProfileRequest, path string, id int) error {
	var sb strings.Builder
	sb.WriteString("UPDATE users SET ")
	args := []any{}

	if path != "" {
		if len(args) > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "avatar = $%d", len(args)+1)
		args = append(args, path)
	}

	if req.Name != "" {
		if len(args) > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "name = $%d", len(args)+1)
		args = append(args, req.Name)
	}

	if req.Bio != "" {
		if len(args) > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "bio = $%d", len(args)+1)
		args = append(args, req.Bio)
	}

	if len(args) == 0 {
		return apperror.ErrNoFieldsToUpdate
	}

	fmt.Fprintf(&sb, " WHERE account_id = $%d", len(args)+1)
	args = append(args, id)

	_, err := db.Exec(ctx, sb.String(), args...)
	if err != nil {
		log.Println("ERROR [repostory:user] failed to update profile:", err)
		return err
	}

	return nil
}

func (ur *UserRepository) GetUsers(ctx context.Context, db DBTX) ([]model.Users, error) {
	sql := "SELECT account_id as id, name FROM users"

	rows, err := db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.Users
	for rows.Next() {
		var user model.Users
		if err := rows.Scan(
			&user.ID,
			&user.Name,
		); err != nil {
			log.Println("ERROR [repostory:user] failed to get users:", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) FollowUser(ctx context.Context, db DBTX, followerID, followedID int) error {
	sql := `
		INSERT INTO
		    followers (following_user_id, followed_user_id)
		VALUES
		    ($1, $2) ON CONFLICT (following_user_id, followed_user_id) DO
		UPDATE
		SET
		    deleted_at = NULL,
		    updated_at = now ();
	`

	_, err := db.Exec(ctx, sql, followerID, followedID)
	if err != nil {
		log.Println("ERROR [repository:follow] failed to follow user:", err)
		return err
	}

	return nil
}
