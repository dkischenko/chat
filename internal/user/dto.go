package user

type UserDTO struct {
	Username string `json:"userName" validate:"required,lowercase,alpha"`
	Password string `json:"password" validate:"required,alphanum"`
}
