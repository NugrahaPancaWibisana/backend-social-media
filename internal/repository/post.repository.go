package repository

import (
	"context"
	"log"

	"github.com/NugrahaPancaWibisana/backend-social-media/internal/dto"
	"github.com/NugrahaPancaWibisana/backend-social-media/internal/model"
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

func (pr *PostRepository) GetFeedPosts(ctx context.Context, db DBTX, userID int) ([]model.FeedPost, error) {
	sql := `
		SELECT
		    p.id,
		    p.content,
		    p.caption,
		    u.account_id,
		    u.name,
		    u.avatar,
		    COUNT(DISTINCT l.user_id) AS likes_count,
		    COUNT(c.user_id) AS comments_count,
		    p.created_at
		FROM
		    posts p
		    JOIN followers f ON f.followed_user_id = p.user_id
		    AND f.following_user_id = $1
		    AND f.deleted_at IS NULL
		    JOIN users u ON u.account_id = p.user_id
		    LEFT JOIN likes l ON l.post_id = p.id
		    AND l.deleted_at IS NULL
		    LEFT JOIN comments c ON c.post_id = p.id
		    AND c.deleted_at IS NULL
		WHERE
		    p.deleted_at IS NULL
		GROUP BY
		    p.id,
		    u.account_id,
		    u.name,
		    u.avatar
		ORDER BY
		    p.created_at DESC;
	`

	rows, err := db.Query(ctx, sql, userID)
	if err != nil {
		log.Println("ERROR [repository:feed] failed to get feed:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.FeedPost
	for rows.Next() {
		var p model.FeedPost
		if err := rows.Scan(
			&p.ID,
			&p.Content,
			&p.Caption,
			&p.UserID,
			&p.UserName,
			&p.UserAvatar,
			&p.LikesCount,
			&p.CommentsCount,
			&p.CreatedAt,
		); err != nil {
			log.Println("ERROR [repository:feed] scan error:", err)
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (pr *PostRepository) CreateLike(ctx context.Context, db DBTX, postID string, userID int) error {
	sql := `
		INSERT INTO 
			likes (post_id, user_id)
		VALUES 
			($1, $2)
	`

	if _, err := db.Exec(ctx, sql, postID, userID); err != nil {
		log.Println("ERROR [repostory:post] failed to like post:", err)
		return err
	}

	return nil
}

func (pr *PostRepository) CreateComment(ctx context.Context, db DBTX, req dto.CreateCommentRequest, userID int) error {
	sql := `
		INSERT INTO 
			comments (comment, post_id, user_id)
		VALUES ($1, $2, $3)
	`

	if _, err := db.Exec(ctx, sql, req.Comment, req.PostID, userID); err != nil {
		log.Println("ERROR [repostory:post] failed to comment post:", err)
		return err
	}

	return nil
}
