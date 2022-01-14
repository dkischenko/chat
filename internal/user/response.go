package user

type UserCreateResponse struct {
	ID       string `json:"id"`
	Username string `json:"userName"`
}

type UserLoginResponse struct {
	Url string `json:"url"`
}
