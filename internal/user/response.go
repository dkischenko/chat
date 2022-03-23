package user

type UserCreateResponse struct {
	ID       string `json:"id"`
	Username string `json:"userName"`
}

type UserLoginResponse struct {
	Url string `json:"url" example:"ws://fancy-chat.io/ws&token=one-time-token"`
}

type UserOnlineResponse struct {
	Count int `json:"count" example:"0"`
}
