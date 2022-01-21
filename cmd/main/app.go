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
	var cfg *config.Config
	var dbHost,
		dbPort,
		dbUsername,
		dbPassword,
		dbDatabase string

	if configPath := os.Getenv("CONFIG"); configPath != "" {
		cfg = config.GetConfig(configPath, &config.Config{})
		dbHost = cfg.Storage.Host
		dbPort = cfg.Storage.Port
		dbUsername = cfg.Storage.Username
		dbPassword = cfg.Storage.Password
		dbDatabase = cfg.Storage.Database
	} else {
		dbHost = os.Getenv("DB_HOST")
		dbPort = os.Getenv("DB_PORT")
		dbUsername = os.Getenv("DB_USERNAME")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbDatabase = os.Getenv("DB_DATABASE")
	}
	l.Entry.Info("Create database connection")
	//storage := database.NewStorage(l)
	//mongoDBCfg := cfg.Storage
	//client, err := mongodb.NewClient(context.Background(), mongoDBCfg.Host, mongoDBCfg.Port, mongoDBCfg.Username,
	//	mongoDBCfg.Password, mongoDBCfg.Database, mongoDBCfg.Options.AuthDB)
	//if err != nil {
	//	panic(err)
	//}
	//storage := database.NewStorage(client, mongoDBCfg.Options.Collection, l)
	client, err := postgres.NewClient(context.Background(), dbHost, dbPort, dbUsername, dbPassword, dbDatabase)
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
