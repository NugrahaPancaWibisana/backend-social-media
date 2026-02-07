package model

import "time"

type FeedPost struct {
	ID            string
	Content       string
	Caption       string
	UserID        int
	UserName      string
	UserAvatar    *string
	LikesCount    int
	CommentsCount int
	CreatedAt     time.Time
}
