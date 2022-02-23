package handlers

import (
	"net/http"
)

//go:generate mockgen -source=handler.go -destination=mocks/handler_mock.go
type Handler interface {
	Register(router *http.ServeMux)
}
