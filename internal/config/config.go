package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string         `yaml:"env" env-default:"development"`
	Database   DatabaseConfig `yaml:"database"`
	HTTPServer HTTPServer     `yaml:"http_server"`
	JWT        JWTConfig      `yaml:"jwt"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	DBName   string `yaml:"dbname" env:"DB_NAME" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
}

type JWTConfig struct {
	Secret   string        `yaml:"secret" env:"JWT_SECRET" env-required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"24h"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}
	return &cfg
}
