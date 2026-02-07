package repository

import (
	"context"
	"log"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
)

type PostRepository struct{}

func NewPostRepository() *PostRepository {
	return &PostRepository{}
}

func (pr *PostRepository) CreatePost(ctx context.Context, db DBTX, req dto.PostRequest, path string, id int) error {
	sql := `
		INSERT INTO
		    posts (content, caption, user_id)
		VALUES
		    ($1, $2, $3)
	`

	if _, err := db.Exec(ctx, sql, path, req.Caption, id); err != nil {
		log.Println("ERROR [repostory:post] failed to create post:", err)
		return err
	}

	return nil
}
