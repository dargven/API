package user

import (
	"API/internal/services/userService"
	"log/slog"
)

type Handler struct { // Вынести в отдельный handler
	logger  *slog.Logger
	Service *userService.UserService
}
