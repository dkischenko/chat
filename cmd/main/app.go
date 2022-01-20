package main

import (
	"context"
	"github.com/dkischenko/chat/internal/app"
	"github.com/dkischenko/chat/internal/config"
	"github.com/dkischenko/chat/internal/user"
	database "github.com/dkischenko/chat/internal/user/database/postgres"
	"github.com/dkischenko/chat/pkg/database/postgres"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"os"
)

func main() {
	l, err := logger.GetLogger()
	if err != nil {
		panic(err)
	}
	l.Entry.Info("Create router")
	router := httprouter.New()
	configPath := os.Getenv("CONFIG")
	cfg := config.GetConfig(configPath, &config.Config{})
	l.Entry.Info("Create database connection")
	//storage := database.NewStorage(l)
	//mongoDBCfg := cfg.Storage
	//client, err := mongodb.NewClient(context.Background(), mongoDBCfg.Host, mongoDBCfg.Port, mongoDBCfg.Username,
	//	mongoDBCfg.Password, mongoDBCfg.Database, mongoDBCfg.Options.AuthDB)
	//if err != nil {
	//	panic(err)
	//}
	//storage := database.NewStorage(client, mongoDBCfg.Options.Collection, l)
	storageCfg := cfg.Storage
	client, err := postgres.NewClient(context.Background(), storageCfg.Host, storageCfg.Port, storageCfg.Username,
		storageCfg.Password, storageCfg.Database)
	if err != nil {
		panic(err)
	}
	storage := database.NewStorage(client, l)
	service := user.NewService(l, storage)
	l.Entry.Info("Register user handler")
	handler := user.NewHandler(l, service)
	handler.Register(router)

	app.Run(router, l, cfg)
}
