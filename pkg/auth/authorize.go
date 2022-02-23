package auth

//go:generate mockgen -source=authorize.go -destination=mocks/authorize_mock.go
type Authorize interface {
	CreateJWT(userId string) (string, error)
	ParseJWT(token string) (string, error)
}
