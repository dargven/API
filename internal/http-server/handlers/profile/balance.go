package profile

import (
	authMiddleware "API/internal/http-server/middleware/auth"
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

// BalanceUpdater интерфейс для пополнения баланса
type BalanceUpdater interface {
	UpdateUserBalance(userID int64, amount float64) error
}

// TopUpBalanceRequest запрос на пополнение баланса
type TopUpBalanceRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// NewTopUpBalance возвращает хендлер для пополнения баланса
// @Summary Пополнить баланс
// @Description Пополняет баланс текущего пользователя
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body TopUpBalanceRequest true "Сумма пополнения"
// @Success 200 {object} resp.Response
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /profile/balance [post]
func NewTopUpBalance(log *slog.Logger, updater BalanceUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.profile.TopUpBalance"

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

		var req TopUpBalanceRequest
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

		if err := updater.UpdateUserBalance(userID, req.Amount); err != nil {
			log.Error("failed to top up balance", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to top up balance"))
			return
		}

		log.Info("balance topped up", slog.Float64("amount", req.Amount), slog.Int64("user_id", userID))
		render.JSON(w, r, resp.OK())
	}
}
