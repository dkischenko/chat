package config

import (
	"fmt"
	v "github.com/dkischenko/chat/internal/validator"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Listen struct {
		Ip   string `yaml:"ip" env-default:"0.0.0.0" validate:"required,ip"`
		Port string `yaml:"port" env-default:"8080" validate:"required,numeric"`
	} `yaml:"listen"`
	Storage struct {
		Host     string `yaml:"host" validate:"required,alpha"`
		Port     string `yaml:"port" validate:"required,numeric"`
		Username string `yaml:"username" validate:"required"`
		Password string `yaml:"password" validate:"required"`
		Database string `yaml:"database" validate:"required"`
		Options  struct {
			Collection string `yaml:"collection" validate:"omitempty,alpha"`
			AuthDB     string `yaml:"auth_db" validate:"omitempty,alpha"`
		} `yaml:"options"`
	} `yaml:"storage"`
	Auth struct {
		AccessTokenTTL string `yaml:"accessTokenTTL" validate:"required"`
	} `yaml:"auth"`
	WS struct {
		WsHost string `yaml:"wsHost" validate:"required"`
	} `yaml:"ws"`
}

func GetConfig(cfgPath string, instance *Config) *Config {
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

	err = validateConfig(instance)
	if err != nil {
		l.Entry.Errorf("config error with: %s", err)
		panic(err)
	}

	return instance
}

func validateConfig(cfg *Config) (err error) {
	valid := v.New()
	err = valid.Vld.Struct(cfg)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			return err
		}

		// from here you can create your own error messages in whatever language you wish

	}

	return
}

// Fill config if config file absent and it fills with env variables.
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
	cfg.WS.WsHost = os.Getenv("WS_HOST")
}
