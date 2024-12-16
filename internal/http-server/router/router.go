package httpserver

import (
	"API/internal/Storage/postrgeSQL"
	"API/internal/config"
	"API/internal/http-server/handlers/booking"
	userHandler "API/internal/http-server/handlers/user"
	"API/repositories/userRepository"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
)

func NewRouter(cfg *config.Config, logger *slog.Logger, db *postrgeSQL.Database) *chi.Mux {
	r := chi.NewRouter()

	// Настраиваем middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// Настройка репозиториев и хендлеров
	userRepo := userRepository.NewUserRepository(db)
	newUserHandler := userHandler.NewUserHandler(userRepo)

	// Маршруты для пользователей
	r.Route("/users", func(r chi.Router) {
		r.Get("/{user_id}", newUserHandler.GetUserByIDHandler)
		r.Post("/", newUserHandler.CreateUserHandler)
		r.Delete("/{user_id}", newUserHandler.DeleteUserHandler)
	})

	// Маршруты для бронирований
	handlerB := booking.Handler{}
	r.Get("/events/{event_id}", handlerB.GetEventByID)
	r.Get("/events/", handlerB.AllEvents)

	// Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Тестовый маршрут
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Welcome to API"))
		if err != nil {
			return
		}
	})

	return r
}
