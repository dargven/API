package logger

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/middleware"
)

//Его вообще не обязательно писать, можем обойтись и стандартным router.Use(middleware.Logger)

func New(log *slog.Logger) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)
		log.Info("logger middleware enabled")

		//Код самого обработчика
		fn := func(w http.ResponseWriter, r *http.Request) {
			//Собираем исходную информацию о запросе
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			// создаем обертку вокруг `http.ResponseWriter`
			// для получения сведений об ответе
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			// Запись отправится в лог в defer
			// в этот момент запрос уже будет обработан
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()
			// Передаем управление следующему обработчику в цепочке middleware
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
