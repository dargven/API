package main

import (
	_ "API/docs"
	"API/internal/Storage/postrgeSQL"
	"API/internal/config"
	"API/internal/http-server/handlers/booking"
	"API/internal/http-server/handlers/test"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"

	"log/slog"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type DataBaseConfig struct {
	Name     string
	Password string
	DBName   string
	Port     string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	checkEnvVars()
	cfg := config.MustLoad()
	logger := setupLogger(cfg.Env)
	logger = logger.With(slog.String("env", cfg.Env)) //к каждому сообщению будет добавляться поле с информацией о текущем окружении
	address := cfg.HTTPServer.Address
	handlerB := booking.Handler{}
	logger.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	logger.Debug("logger debug mode enabled")
	db, err := postrgeSQL.NewDatabase(&cfg.DataBase)
	if err != nil {
		logger.Error("Failed to connect to database: %v", err)
	}
	if db != nil {
		defer db.Close()
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга(понимания сколько выполняется каждый запрос)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	r.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	r.Route("/book", func(r chi.Router) {
		r.Use(middleware.BasicAuth("API", map[string]string{
			cfg.DataBase.User: cfg.DataBase.Password,
		}))

		// сюда добавить Post соответственно
		//насрать хендлерами для удаления и т.д.
	}) //если не работает до вместо database ебануть HTTPServer

	r.Get("/", test.GetRootHandler)
	r.Get("/book/{event_id}", handlerB.GetEventByID)
	r.Get("/book/events", handlerB.AllEvents)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	fmt.Printf("Started server at %s\n", cfg.HTTPServer.Address)
	if err := http.ListenAndServe(address, r); err != nil {
		logger.Error("Error starting server: %s", err)
	}
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
	}
	return logger
}

func checkEnvVars() {
	requiredVars := []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_PORT", "POSTGRES_NAME", "ENV"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Environment variable %s is not set", v)
		}
	}
}
