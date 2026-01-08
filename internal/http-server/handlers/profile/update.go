package profile

import (
	authMiddleware "API/internal/http-server/middleware/auth"
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"API/internal/models"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

// ProfileUpdater интерфейс для обновления профиля
type ProfileUpdater interface {
	UpdateUserProfile(userID int64, name string, phone, avatarURL, bio *string) (*models.User, error)
}

// UpdateProfileRequest запрос на обновление профиля
type UpdateProfileRequest struct {
	Name      string  `json:"name" validate:"required,min=2,max=100"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
	Bio       *string `json:"bio,omitempty" validate:"omitempty,max=500"`
}

// UpdateProfileResponse ответ на обновление профиля
type UpdateProfileResponse struct {
	resp.Response
	Profile models.ProfileResponse `json:"profile"`
}

// NewUpdate возвращает хендлер для обновления профиля
// @Summary Обновить профиль
// @Description Обновляет профиль текущего пользователя
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body UpdateProfileRequest true "Данные профиля"
// @Success 200 {object} UpdateProfileResponse
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /profile [put]
func NewUpdate(log *slog.Logger, updater ProfileUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.profile.Update"

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

		var req UpdateProfileRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request body"))
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("validation failed", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		user, err := updater.UpdateUserProfile(userID, req.Name, req.Phone, req.AvatarURL, req.Bio)
		if err != nil {
			log.Error("failed to update profile", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update profile"))
			return
		}

		render.JSON(w, r, UpdateProfileResponse{
			Response: resp.OK(),
			Profile:  user.ToProfileResponse(),
		})
	}
}
