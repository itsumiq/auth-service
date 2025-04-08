package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	Database database
}

type database struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     int    `env:"DB_PORT" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Name     string `env:"DB_DATABASE" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
}

var (
	cfg  config
	once sync.Once
)

func loadConfig() {
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(".env file not loaded")
	}
	fmt.Println("Config is loaded")
}

func Get() *config {
	once.Do(loadConfig)
	return &cfg
}
