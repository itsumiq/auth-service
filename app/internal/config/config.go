package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database database
	Token    token
	Server   server
}

type database struct {
	Host     string `env:"DB_HOST"     env-required:"true"`
	Port     int    `env:"DB_PORT"     env-required:"true"`
	User     string `env:"DB_USER"     env-required:"true"`
	Name     string `env:"DB_DATABASE" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
}

type token struct {
	SecretKey string `env:"TOKEN_SECRET_KEY" env-required:"true"`
}

type server struct {
	TimeoutResponse uint `env:"SERVER_TIMEOUT_RESPONSE" env-required:"true"`
}

var (
	cfg  Config
	once sync.Once
)

func loadConfig() {
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(".env file not loaded")
	}
	fmt.Println("Config is loaded")
}

func Get() *Config {
	once.Do(loadConfig)
	return &cfg
}
