package dto

type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email" example:"user@example.com"`
	Password        string `json:"password" binding:"required,min=8" example:"user123@"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password" example:"user123@"`
}
