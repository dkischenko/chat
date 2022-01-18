package config

import (
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Listen struct {
		Ip   string `yaml:"ip" env-default:"127.0.0.1"`
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
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		l := logger.GetLogger()
		l.Info("Start read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, errGD := cleanenv.GetDescription(instance, nil)
			if errGD != nil {
				l.Fatalf("GetDescription error: %s", errGD)
			}
			l.Info(help)
			l.Fatal(err)
		}
	})

	return instance
}
