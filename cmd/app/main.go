package main

import (
	"API/internal/Storage/postrgeSQL"
	"API/internal/config"
	"API/internal/http-server/handlers/test"
	"fmt"
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
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env)) //к каждому сообщению будет добавляться поле с информацией о текущем окружении
	address := cfg.HTTPServer.Address
	log.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")
	db, err := postrgeSQL.NewDatabase(&cfg.DataBase)
	if err != nil {
		log.Error("Failed to connect to database: %v", err)
	}
	if db != nil {
		defer db.Close()
	}
	//}

	router := chi.NewRouter()
	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга(понимания сколько выполняется каждый запрос)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received /test request")
		w.Write([]byte("Server is working"))
	})

	router.Get("/hello", test.HelloHandler(log))

	fmt.Printf("Starting server at %s\n", cfg.HTTPServer.Address)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Error("Error starting server: %s", err)
	}
	// router.Post("/", save.New(log, db))
	// // Все пути этого роутера будут начинаться с префикса `/url`
	// router.Route("/url", func(r chi.Router) {
	// 	// Подключаем авторизацию
	// 	r.Use(middleware.BasicAuth("API", map[string]string{
	// 		// Передаем в middleware креды
	// 		cfg.HTTPServer.User: cfg.HTTPServer.Password,
	// 		// Если у вас более одного пользователя,
	// 		// то можете добавить остальные пары по аналогии.
	// 	}))

	// 	r.Post("/", save.New(log, db))
	// })

	// Хэндлер redirect остается снаружи, в основном роутере
	//router.Get("/{alias}", redirect.New(log, db))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

// func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		const op = "handlers.url.redicted.New"

// 		log = log.With(
// 			slog.String("op", op),
// 			slog.String("request_id", middleware.GetReqID(r.Context())),
// 		)

// 		// Роутер chi позволяет делать вот такие финты -
// 		// получать GET-параметры по их именам.
// 		// Имена определяются при добавлении хэндлера в роутер, это будет ниже.
// 		alias := chi.URLParam(r, "alias")
// 		if alias == "" {
// 			log.Error("alias is empty")

// 			render.JSON(w, r, resp.Error("alias is empty"))

// 			return
// 		}

// 		// Находим URL по алиасу в БД
// 		resURL, err := urlGetter.GetURL(alias)
// 		if errors.Is(err, storage.ErrURLNotFound) {
// 			// Не нашли URL, сообщаем об этом клиенту
// 			log.Error("request body is empty")

// 			render.JSON(w, r, resp.Error("empty request"))

// 			return

// 		}
// 		if err != nil {
// 			log.Error("failed to get URL", sl.Err(err))

// 			render.JSON(w, r, resp.Error("internal error"))

// 			return

// 		}

// 		log.Info("got url", slog.String("url", resURL))

// 		http.Redirect(w, r, resURL, http.StatusFound)
// 	}
// }

func checkEnvVars() {
	requiredVars := []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_PORT", "POSTGRES_NAME", "ENV"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Environment variable %s is not set", v)
		}
	}
}
