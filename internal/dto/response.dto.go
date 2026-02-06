package dto

type Response struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Data retrieved successfully"`
}

type ResponseSuccess struct {
	Response
	Data any `json:"data"`
}

type ResponseSuccessWithMeta struct {
	Response
	Data any `json:"data"`
	Meta any `json:"meta"`
}

type ResponseError struct {
	Response
	Error string `json:"error" example:"Internal Server Error"`
}
