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

type RequestVerificationEmailDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type TokenDTO struct {
	Token string `json:"token"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RequestChangePasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type ChangePasswordDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type CreatePersonDto struct {
	Name   string `json:"name" validate:"required"`
	Avatar string `json:"avatar"`
}

type CreateTagDto struct {
	Title        string `json:"title" validate:"required"`
	IsForNote    bool   `json:"isForNote"`
	IsForMeeting bool   `json:"isForMeeting"`
}
