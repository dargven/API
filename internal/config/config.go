package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-default:"development"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	DataBase    `yaml:"data_base"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"` //чет адресс подозрительный
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type DataBase struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"PORT"`
	User     string `yaml:"user" env-default:"USER"`
	Password string `yaml:"password" env-default:"PASSWORD"` // В иделе password хранить в .env. В Prod придется переделывать
	Name     string `yaml:"name" env-default:"NAME"`
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

//func (cfg *DataBase) BuildPostgresDSN() string {
//	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
//		cfg.DBUser,
//		cfg.DBPassword,
//		cfg.DBHost,
//		cfg.DBPort,
//		cfg.DBName,
//	)
//}
