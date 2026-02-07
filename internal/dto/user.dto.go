package dto

type User struct {
	ID     int    `json:"id" example:"1"`
	Email  string `json:"email" example:"user@example.com"`
	Name   string `json:"name" example:"Jhon Doe"`
	Avatar string `json:"avatar" example:"/image.png"`
	Bio    string `json:"bio" example:"my bio"`
}

type Users struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Jhon Doe"`
}
