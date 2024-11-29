package userHandler

import (
	"API/internal/services/userService"
	"log/slog"
)

type UserHandler struct { // Вынести в отдельный handler
	logger  *slog.Logger
	Service *userService.UserService
}
