package main

import (
	"context"
	"github.com/dkischenko/chat/internal/app"
	"github.com/dkischenko/chat/internal/config"
	"github.com/dkischenko/chat/internal/user"
	database "github.com/dkischenko/chat/internal/user/database/postgres"
	"github.com/dkischenko/chat/pkg/database/postgres"
	"github.com/dkischenko/chat/pkg/logger"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/swaggo/swag/example/celler/docs"
)

// @title            Fancy Golang chat
// @version          1.0.0
// @description      Just a simple chat service
// @license.name     Apache 2.0
// @license.url      http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath         /v1
// @tag.name         user
// @tag.description  Operations about user
// @tag.name         chat
// @tag.description  Operations about chat
func main() {
	l, err := logger.GetLogger()
	if err != nil {
		panic(err)
	}
	l.Entry.Info("Create router")
	router := http.NewServeMux()
	var cfg *config.Config

	configPath := os.Getenv("CONFIG")
	var once sync.Once
	once.Do(func() {
		cfg = config.GetConfig(configPath, &config.Config{})
	})

	l.Entry.Info("Create database connection")

	client, err := postgres.NewClient(context.Background(), cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username,
		cfg.Storage.Password, cfg.Storage.Database)
	if err != nil {
		panic(err)
	}
	storage := database.NewStorage(client, l)
	accessTokenTTL, err := time.ParseDuration(cfg.Auth.AccessTokenTTL)
	if err != nil {
		panic(err)
	}
	s := user.NewService(l, storage, accessTokenTTL)
	l.Entry.Info("Register user handler")
	handler := user.NewHandler(l, s, cfg)
	handler.Register(router)

	app.Run(router, l, cfg)
}
