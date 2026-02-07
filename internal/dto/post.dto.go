package dto

import "time"

type FeedPost struct {
	ID            string    `json:"id"`
	Content       string    `json:"content"`
	Caption       string    `json:"caption"`
	Likes         int       `json:"likes"`
	Comments      int       `json:"comments"`
	Author        Author    `json:"author"`
	CreatedAt     time.Time `json:"created_at"`
}

type Author struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar,omitempty"`
}
