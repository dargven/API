package profile

import (
	authMiddleware "API/internal/http-server/middleware/auth"
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"API/internal/models"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

// ProfileGetter интерфейс для получения профиля
type ProfileGetter interface {
	GetUserByID(id int64) (*models.User, error)
}

// GetProfileResponse ответ с профилем
type GetProfileResponse struct {
	resp.Response
	Profile models.ProfileResponse `json:"profile"`
}

// NewGet возвращает хендлер для получения профиля
// @Summary Получить профиль
// @Description Возвращает профиль текущего пользователя
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} GetProfileResponse
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /profile [get]
func NewGet(log *slog.Logger, getter ProfileGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.profile.Get"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID, ok := authMiddleware.GetUserIDFromContext(r.Context())
		if !ok {
			log.Error("user_id not found in context")
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("unauthorized"))
			return
		}

		user, err := getter.GetUserByID(userID)
		if err != nil {
			log.Error("failed to get user", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get profile"))
			return
		}

		render.JSON(w, r, GetProfileResponse{
			Response: resp.OK(),
			Profile:  user.ToProfileResponse(),
		})
	}
}
