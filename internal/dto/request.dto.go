package dto

import "mime/multipart"

type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email" example:"user@example.com"`
	Password        string `json:"password" binding:"required,min=8" example:"user123@"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,eqfield=Password" example:"user123@"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"user123@"`
}

type UpdateProfileRequest struct {
	Avatar *multipart.FileHeader `form:"avatar"`
	Name   string                `form:"name" binding:"omitempty,min=3" example:"John Doe"`
	Bio    string                `form:"bio" binding:"omitempty,min=3" example:"my bio"`
}

type PostRequest struct {
	Content *multipart.FileHeader `form:"content"`
	Caption string                `form:"caption" binding:"omitempty,min=3" example:"my caption"`
}

type LikeRequest struct {
	PostID string `json:"post_id" binding:"required"`
}

type CreateCommentRequest struct {
	PostID  string    `json:"post_id" binding:"required"`
	Comment string `json:"comment" binding:"required,min=1"`
}