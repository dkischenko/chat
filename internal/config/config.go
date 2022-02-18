package config

import (
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"sync"
)

type Config struct {
	Listen struct {
		Ip   string `yaml:"ip" env-default:"0.0.0.0"`
		Port string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	Storage struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Options  struct {
			Collection string `yaml:"collection"`
			AuthDB     string `yaml:"auth_db"`
		} `yaml:"options"`
	} `yaml:"storage"`
	Auth struct {
		AccessTokenTTL string `yaml:"accessTokenTTL"`
	} `yaml:"auth"`
}

func GetConfig(cfgPath string, instance *Config) *Config {
	var once sync.Once
	once.Do(func() {
		l, err := logger.GetLogger()
		if err != nil {
			panic("Error create logger")
		}
		l.Entry.Info("Start read application config")
		instance = &Config{}
		if cfgPath != "" {
			if err := cleanenv.ReadConfig(cfgPath, instance); err != nil {
				help, errGD := cleanenv.GetDescription(instance, nil)
				if errGD != nil {
					l.Entry.Fatalf("GetDescription error: %s", errGD)
				}
				l.Entry.Info(help)
				l.Entry.Fatal(err)
			}
		} else {
			populateConfig(instance)
		}
	})

	return instance
}

func populateConfig(cfg *Config) {
	cfg.Storage.Host = os.Getenv("DB_HOST")
	cfg.Storage.Port = os.Getenv("DB_PORT")
	cfg.Storage.Username = os.Getenv("DB_USERNAME")
	cfg.Storage.Password = os.Getenv("DB_PASSWORD")
	cfg.Storage.Database = os.Getenv("DB_DATABASE")
	cfg.Storage.Options.AuthDB = os.Getenv("DB_AUTHDB")
	cfg.Storage.Options.Collection = os.Getenv("DB_COLLECTION")
	cfg.Listen.Ip = os.Getenv("APP_IP")
	cfg.Listen.Port = os.Getenv("PORT")
	cfg.Auth.AccessTokenTTL = os.Getenv("ACCESSTOKENTTL")
}
