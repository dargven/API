package main

import (
	_ "API/docs"
	"API/internal/config"
	httpserver "API/internal/http-server/router"
	"API/internal/storage/postrgesql"
	"errors"
	"log"
	"net/http"
	"os"

	"log/slog"

	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	checkRequiredEnvVars()

	// Загрузка конфигурации
	cfg := config.MustLoad()

	// Настройка логгера
	logger := setupLogger(cfg.Env)
	logger = logger.With(slog.String("env", cfg.Env))

	logger.Info("Starting application", slog.String("environment", cfg.Env), slog.String("address", cfg.HTTPServer.Address))

	// Настройка подключения к базе данных
	db, err := postrgesql.NewDatabase(&cfg.DataBase)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Инициализация роутера
	router := httpserver.NewRouter(cfg, logger, db)

	// Запуск сервера
	startServer(cfg, router, logger)
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log.Fatalf("Invalid environment: %s", env)
	}
	return logger
}

func checkRequiredEnvVars() {
	requiredVars := []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_PORT", "POSTGRES_NAME", "ENV"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Environment variable %s is not set", v)
		}
	}
}

func startServer(cfg *config.Config, router http.Handler, logger *slog.Logger) {
	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	logger.Info("HTTP server is starting", slog.String("address", cfg.HTTPServer.Address))

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Failed to start HTTP server", "error", err)
		os.Exit(1)
	}
}
