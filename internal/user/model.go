package user

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Username     string `json:"userName" bson:"username"`
	PasswordHash string `json:"passwordHash" bson:"password"`
}

type UserDTO struct {
	Username string `json:"userName" validate:"required,lowercase,alpha"`
	Password string `json:"password" validate:"required,alphanum"`
}
