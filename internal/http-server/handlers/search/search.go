package search

import (
	resp "API/internal/lib/api/response"
	"API/internal/lib/logger/sl"
	"API/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

// EventSearcher интерфейс для поиска мероприятий
type EventSearcher interface {
	SearchEvents(query string, category string, dateFrom, dateTo *time.Time, priceMin, priceMax *float64, limit, offset int) ([]*models.Event, int, error)
}

// SearchResponse ответ с результатами поиска
type SearchResponse struct {
	resp.Response
	Events  []models.EventResponse `json:"events"`
	Total   int                    `json:"total"`
	Limit   int                    `json:"limit"`
	Offset  int                    `json:"offset"`
	HasMore bool                   `json:"has_more"`
}

// NewSearch возвращает хендлер для поиска мероприятий
// @Summary Поиск мероприятий
// @Description Полнотекстовый поиск мероприятий с фильтрами
// @Tags search
// @Security BearerAuth
// @Produce json
// @Param q query string false "Поисковый запрос"
// @Param category query string false "Категория (concert, sport, theater, exhibition, festival, other)"
// @Param date_from query string false "Дата от (RFC3339)"
// @Param date_to query string false "Дата до (RFC3339)"
// @Param price_min query number false "Минимальная цена"
// @Param price_max query number false "Максимальная цена"
// @Param limit query int false "Лимит (default 20, max 100)"
// @Param offset query int false "Смещение (default 0)"
// @Success 200 {object} SearchResponse
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /search [get]
func NewSearch(log *slog.Logger, searcher EventSearcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.search.Search"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Парсинг параметров
		query := r.URL.Query().Get("q")
		category := r.URL.Query().Get("category")

		// Парсинг дат
		var dateFrom, dateTo *time.Time
		if df := r.URL.Query().Get("date_from"); df != "" {
			t, err := time.Parse(time.RFC3339, df)
			if err != nil {
				log.Error("invalid date_from format", sl.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("invalid date_from format, use RFC3339"))
				return
			}
			dateFrom = &t
		}
		if dt := r.URL.Query().Get("date_to"); dt != "" {
			t, err := time.Parse(time.RFC3339, dt)
			if err != nil {
				log.Error("invalid date_to format", sl.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("invalid date_to format, use RFC3339"))
				return
			}
			dateTo = &t
		}

		// Парсинг цен
		var priceMin, priceMax *float64
		if pm := r.URL.Query().Get("price_min"); pm != "" {
			p, err := strconv.ParseFloat(pm, 64)
			if err != nil || p < 0 {
				log.Error("invalid price_min", sl.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("invalid price_min"))
				return
			}
			priceMin = &p
		}
		if pm := r.URL.Query().Get("price_max"); pm != "" {
			p, err := strconv.ParseFloat(pm, 64)
			if err != nil || p < 0 {
				log.Error("invalid price_max", sl.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("invalid price_max"))
				return
			}
			priceMax = &p
		}

		// Пагинация
		limit := 20
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
				limit = parsed
			}
		}

		offset := 0
		if o := r.URL.Query().Get("offset"); o != "" {
			if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
				offset = parsed
			}
		}

		events, total, err := searcher.SearchEvents(query, category, dateFrom, dateTo, priceMin, priceMax, limit, offset)
		if err != nil {
			log.Error("failed to search events", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to search events"))
			return
		}

		response := make([]models.EventResponse, 0, len(events))
		for _, e := range events {
			response = append(response, e.ToResponse())
		}

		log.Info("search completed",
			slog.String("query", query),
			slog.Int("total", total),
			slog.Int("returned", len(response)),
		)

		render.JSON(w, r, SearchResponse{
			Response: resp.OK(),
			Events:   response,
			Total:    total,
			Limit:    limit,
			Offset:   offset,
			HasMore:  offset+len(response) < total,
		})
	}
}
