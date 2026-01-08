package main

import (
	storage "API/internal/Storage"
	"API/internal/Storage/postgres"
	"API/internal/auth"
	"API/internal/config"
	authHandlers "API/internal/http-server/handlers/auth"
	"API/internal/http-server/handlers/events"
	"API/internal/http-server/handlers/url/save"
	authMiddleware "API/internal/http-server/middleware/auth"
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"errors"
	"net/http"
	"os"

	_ "API/docs"

	"golang.org/x/exp/slog"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// URLGetter is an interface for getting url by alias.
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
	log.Info("initializing server", slog.String("address", cfg.HTTPServer.Address))
	log.Debug("logger debug mode enabled")

	storage, err := postgres.New(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}
	defer storage.Close()

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.TokenTTL)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.RedirectSlashes)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandlers.NewRegister(log, storage, jwtManager))
		r.Post("/login", authHandlers.NewLogin(log, storage, jwtManager))
	})

	router.Route("/events", func(r chi.Router) {
		r.Use(authMiddleware.JWTAuth(log, jwtManager))
		r.Post("/", events.NewCreate(log, storage))
		r.Get("/", events.NewGetAll(log, storage))
		r.Get("/{id}", events.NewGetByID(log, storage))
	})

	router.Route("/url", func(r chi.Router) {
		r.Use(authMiddleware.JWTAuth(log, jwtManager))
		r.Post("/", save.New(log, storage))
	})

	// Публичный роут для редиректа
	router.Get("/{alias}", redirectHandler(log, storage))

	// Запуск сервера
	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))
	log.Info("swagger UI available at", slog.String("url", "http://"+cfg.HTTPServer.Address+"/swagger/index.html"))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
		os.Exit(1)
	}
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
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}

// redirectHandler обрабатывает редирект по алиасу
func redirectHandler(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")
			render.JSON(w, r, resp.Error("alias is empty"))
			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Error("url not found", slog.String("alias", alias))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to get URL", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("redirecting", slog.String("alias", alias), slog.String("url", resURL))
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
