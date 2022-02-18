package auth

type Authorize interface {
	CreateJWT(userId string) (string, error)
	ParseJWT(token string) (string, error)
}
