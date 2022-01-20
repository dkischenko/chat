package user

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
