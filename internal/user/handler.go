package user

import (
	"encoding/json"
	"fmt"
	"github.com/dkischenko/chat/internal/config"
	"github.com/dkischenko/chat/internal/handlers"
	"github.com/dkischenko/chat/internal/middleware"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/go-playground/validator/v10"
	"net/http"
	"os"
	"time"
)

const (
	userUrl                = "/v1/user"
	userLoginUrl           = "/v1/user/login"
	userActive             = "/v1/user/active"
	chatUrl                = "/v1/chat/ws.rtm.start/"
	xRateLimit             = "50"
	headerContentType      = "Content-Type"
	headerValueContentType = "application/json"
	headerValueXRateLimit  = "X-Rate-Limit"
	headerXExpiresAfter    = "X-Expires-After"
	wsHostLocal            = "ws://localhost:1000"
	wsHostProd             = "ws://floating-ridge-89522.herokuapp.com"
)

type handler struct {
	logger  *logger.Logger
	service *service
}

func NewHandler(logger *logger.Logger, service *service) handlers.Handler {
	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h *handler) Register(router *http.ServeMux) {
	createUserHandler := http.HandlerFunc(h.CreateUser)
	loginUserHandler := http.HandlerFunc(h.LoginUser)
	activeUserHandler := http.HandlerFunc(h.ActiveUser)
	chatStartHandler := http.HandlerFunc(h.ChatStart)
	router.Handle(userUrl, middleware.PanicAndRecover(middleware.Logging(createUserHandler, h.logger), h.logger))
	router.Handle(userLoginUrl, middleware.PanicAndRecover(middleware.Logging(loginUserHandler, h.logger), h.logger))
	router.Handle(userActive, middleware.PanicAndRecover(middleware.Logging(activeUserHandler, h.logger), h.logger))
	router.Handle(chatUrl, middleware.PanicAndRecover(middleware.Logging(chatStartHandler, h.logger), h.logger))
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	code := h.isPost(r)
	if code > 0 {
		w.WriteHeader(code)
		return
	}
	// @todo: refactor validation to service
	uDTO := &UserDTO{}
	err := json.NewDecoder(r.Body).Decode(uDTO)

	if err != nil {
		h.logger.Entry.Error("wrong json format")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	v := validator.New()

	if err := v.Struct(uDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseBody := ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("got wrong user data: %+v", err),
		}
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			h.logger.Entry.Errorf("problems with encoding data: %+v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		h.logger.Entry.Errorf("got wrong user data: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// @todo: end

	uID, err := h.service.Create(r.Context(), *uDTO)
	if err != nil {
		h.logger.Entry.Errorf("can't create user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// @todo refactor to service
	w.Header().Add(headerContentType, headerValueContentType)
	w.WriteHeader(http.StatusOK)
	responseBody := UserCreateResponse{
		ID:       uID,
		Username: uDTO.Username,
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		h.logger.Entry.Errorf("can't create user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// @todo end

	h.logger.Entry.Infof("create user %+v", uDTO)
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	code := h.isPost(r)
	if code > 0 {
		w.WriteHeader(code)
		return
	}
	// @todo: refactor validation to service
	uDTO := &UserDTO{}
	err := json.NewDecoder(r.Body).Decode(uDTO)

	if err != nil {
		h.logger.Entry.Error("wrong json format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v := validator.New()

	if err := v.Struct(uDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseBody := ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("got wrong user data: %+v", err),
		}
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			h.logger.Entry.Errorf("problems with encoding data: %+v", err)
		}
		h.logger.Entry.Errorf("got wrong user data: %+v", err)
		return
	}
	// @todo: end

	//find user and create token
	u, err := h.service.Login(r.Context(), uDTO)
	if err != nil {
		h.logger.Entry.Errorf("error with user login: %v", err)
	}
	hash, err := h.service.CreateToken(r.Context(), u)
	if err != nil {
		h.logger.Entry.Errorf("error with create token: %v", err)
	}
	// @todo refactor to service
	w.Header().Add(headerValueXRateLimit, xRateLimit)
	configPath := os.Getenv("CONFIG")
	cfg := config.GetConfig(configPath, &config.Config{})
	accessTokenTTL, err := time.ParseDuration(cfg.Auth.AccessTokenTTL)
	if err != nil {
		h.logger.Entry.Errorf("Error with access token ttl: %s", err)
	}

	w.Header().Add(headerXExpiresAfter, time.Now().Local().Add(accessTokenTTL).String())
	w.Header().Add(headerContentType, headerValueContentType)
	w.WriteHeader(http.StatusOK)
	var wsHost string
	if os.Getenv("APP_ENV") == "local" {
		wsHost = wsHostLocal
	} else if os.Getenv("APP_ENV") == "prod" {
		wsHost = wsHostProd
	}
	responseBody := UserLoginResponse{
		Url: wsHost + chatUrl + "?token=" + hash,
	}
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		h.logger.Entry.Errorf("Failed to login user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// @todo end

	h.logger.Entry.Infof("user sussesfully logged in")
}

func (h *handler) ActiveUser(w http.ResponseWriter, r *http.Request) {
	code := h.isGet(r)
	if code > 0 {
		w.WriteHeader(code)
		return
	}

	count, err := h.service.GetOnlineUsers(r.Context())
	if err != nil {
		h.logger.Entry.Error("Error with getting online users count: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseBody := UserOnlineResponse{
		Count: count,
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		h.logger.Entry.Errorf("Failed to login user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) ChatStart(w http.ResponseWriter, r *http.Request) {
	code := h.isGet(r)
	if code > 0 {
		w.WriteHeader(code)
		return
	}

	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) < 1 {
		h.logger.Entry.Error("Url Param 'token' is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// find user with token

	uuid, err := h.service.ParseToken(token[0])
	if err != nil {
		h.logger.Entry.Errorf("please, use valid token: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.service.FindByUUID(r.Context(), uuid)
	if err != nil {
		h.logger.Entry.Errorf("Please, use valid token: %s.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if u.Key == "" {
		h.logger.Entry.Error("May be, your token has been revoked earlier. Please, login again.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ok = h.service.RevokeToken(r.Context(), u)
	if !ok {
		h.logger.Entry.Errorf("failed revoke token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.service.StartWS(w, r, u)
	if err != nil {
		h.logger.Entry.Errorf("wrong http method due error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) isGet(r *http.Request) int {
	if r.Method != "GET" {
		h.logger.Entry.Error("Wrong http method. Use `GET`")
		return http.StatusInternalServerError
	}
	return 0
}

func (h *handler) isPost(r *http.Request) int {
	if r.Method != "POST" {
		h.logger.Entry.Error("Wrong http method. Use `POST`")
		return http.StatusInternalServerError
	}
	return 0
}
