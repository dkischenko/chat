package models

type UserDTO struct {
	Username string `json:"userName" validate:"required,lowercase,alpha" example:"userName"`
	Password string `json:"password" validate:"required,alphanum" example:"password"`
}
