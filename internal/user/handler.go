package user

import (
	"encoding/json"
	"fmt"
	"github.com/dkischenko/chat/internal/handlers"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

const (
	userUrl      = "/v1/user"
	userLoginUrl = "/v1/user/login"
	xRateLimit   = "50"
)

type handler struct {
	handlers.Handler
	logger  *logger.Logger
	service *service
}

func NewHandler(logger *logger.Logger, service *service) handlers.Handler {
	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.POST(userUrl, h.CreateUser)
	router.POST(userLoginUrl, h.LoginUser)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// @todo: refactor validation to service
	uDTO := &UserDTO{}
	err := json.NewDecoder(r.Body).Decode(uDTO)

	if err != nil {
		h.logger.Error("wrong json format")
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
			h.logger.Errorf("problems with encoding data: %+v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		h.logger.Errorf("got wrong user data: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// @todo: end

	uID, err := h.service.Create(r.Context(), *uDTO)
	if err != nil {
		h.logger.Errorf("can't create user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// @todo refactor to service
	w.Header().Add("Content-Type", "application/json")
	// @todo refactor headers
	w.Header().Add("X-Rate-Limit", xRateLimit)
	w.Header().Add("X-Expires-After", time.Now().Local().Add(time.Minute*time.Duration(30)).String())
	w.WriteHeader(http.StatusOK)
	responseBody := UserCreateResponse{
		ID:       uID,
		Username: uDTO.Username,
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		h.logger.Errorf("can't create user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// @todo end

	h.logger.Infof("create user %+v", uDTO)
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// validate income data
	// @todo: refactor validation to service
	uDTO := &UserDTO{}
	err := json.NewDecoder(r.Body).Decode(uDTO)

	if err != nil {
		h.logger.Error("wrong json format")
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
			h.logger.Errorf("problems with encoding data: %+v", err)
		}
		h.logger.Errorf("got wrong user data: %+v", err)
		return
	}
	// @todo: end

	//find user and create token
	hash, err := h.service.Login(r.Context(), uDTO.Username)
	if err != nil {
		h.logger.Errorf("error with user login: %v", err)
	}

	// @todo refactor to service
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseBody := UserLoginResponse{
		Url: "ws://fancy-chat.io/ws&token=" + hash,
	}
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		h.logger.Errorf("Failed to login user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// @todo end

	h.logger.Infof("user sussesfully logged in")
}
