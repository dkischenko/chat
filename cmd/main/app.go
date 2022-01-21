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

	if configPath := os.Getenv("CONFIG"); configPath == "" {
		cfg = &config.Config{}
		cfg.Storage.Host = os.Getenv("DB_HOST")
		cfg.Storage.Port = os.Getenv("DB_PORT")
		cfg.Storage.Username = os.Getenv("DB_USERNAME")
		cfg.Storage.Password = os.Getenv("DB_PASSWORD")
		cfg.Storage.Database = os.Getenv("DB_DATABASE")
		cfg.Storage.Options.AuthDB = os.Getenv("DB_AUTHDB")
		cfg.Storage.Options.Collection = os.Getenv("DB_COLLECTION")
		cfg.Listen.Ip = os.Getenv("APP_IP")
		cfg.Listen.Port = os.Getenv("PORT")
	} else {
		cfg = config.GetConfig(configPath, &config.Config{})
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
	client, err := postgres.NewClient(context.Background(), cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username,
		cfg.Storage.Password, cfg.Storage.Database)
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
