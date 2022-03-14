package uerrors

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrGetOnlineUsers        = errors.New("error with getting online users")
	ErrRevokeToken           = errors.New("error with revoke token")
	ErrEmptyUsername         = errors.New("username can't be empty")
	ErrCreateUser            = errors.New("error with creating user due a database issue")
	ErrFindOneUser           = errors.New("error with finding user")
	ErrCheckUserPasswordHash = errors.New("error with using wrong password")
	ErrFindUserByUIID        = errors.New("error with finding user by UIID")
	ErrCreateJWTToken        = errors.New("error with creation of JWT token of user")
	ErrEmptyUserKey          = errors.New("error with user empty key")
	ErrUserUpdateKey         = errors.New("error with user's updating key")
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
