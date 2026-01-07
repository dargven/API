package httpserver

import (
	"API/internal/config"
	"API/internal/http-server/handlers/booking"
	userHandler "API/internal/http-server/handlers/user"
	service "API/internal/services/userService"
	"API/internal/storage/postrgesql"
	"API/repositories/userRepository"

	"log/slog"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(cfg *config.Config, logger *slog.Logger, db *postrgesql.Database) *chi.Mux {
	r := chi.NewRouter()

	// Настраиваем middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// Настройка репозиториев и хендлеров
	userRepo := userRepository.NewUserRepository(db)
	handler := userHandler.NewUserHandler(userRepo)
	service.NewUserService(userRepo)
	//
	// Маршруты для пользователей
	r.Route("/users", func(r chi.Router) {
		r.Get("/{user_id}", handler.GetUserByIDHandler)
		r.Post("/", handler.CreateUserHandler)
		r.Post("/login", handler.LoginHandler)
		r.Delete("/{user_id}", handler.DeleteUserHandler)
	})

	// Маршруты для бронирований
	handlerB := booking.Handler{}
	r.Get("/events/{event_id}", handlerB.GetEventByID)
	r.Get("/events/", handlerB.AllEvents)

	// Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
