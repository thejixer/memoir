package models

type ResponseDTO struct {
	Msg        string `json:"msg"`
	StatusCode int    `json:"statusCode"`
}

type SignUpDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
