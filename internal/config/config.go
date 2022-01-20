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

func GetConfig(cfgPath string, instance *Config) *Config {
	var once sync.Once
	once.Do(func() {
		l, err := logger.GetLogger()
		if err != nil {
			panic("Error create logger")
		}
		l.Entry.Info("Start read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig(cfgPath, instance); err != nil {
			help, errGD := cleanenv.GetDescription(instance, nil)
			if errGD != nil {
				l.Entry.Fatalf("GetDescription error: %s", errGD)
			}
			l.Entry.Info(help)
			l.Entry.Fatal(err)
		}
	})

	return instance
}
