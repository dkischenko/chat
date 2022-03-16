package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

type Manager struct {
	signinKey []byte
	tokenTTL  time.Duration
}

func NewManager(tokenTTL time.Duration) (*Manager, error) {
	var key string
	if key = os.Getenv("SIGNINKEY"); key == "" {
		return nil, errors.New("empty signin key passed")
	}

	return &Manager{signinKey: []byte(key), tokenTTL: tokenTTL}, nil
}

func (m *Manager) CreateJWT(userId string) (string, error) {
	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(m.tokenTTL).Unix()
	claims["iss_at"] = time.Now().Unix()
	claims["user_id"] = userId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.signinKey)
}

func (m *Manager) ParseJWT(tokenString string) (user_id string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return m.signinKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims["user_id"] == nil {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["user_id"].(string), nil
}
