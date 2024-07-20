package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(slog.String("component", "middleware/logger"))

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			// эта часть выполняется до обработки запроса
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			// создаем обертку вокруг http.ResponseWriter для получения сведений об ответе
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Момент получения запроса, чтобы вычислить время обработки
			timeNow := time.Now()

			// эта часть будет выполнена после окончательной обработки запроса
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(timeNow).String()),
				)
			}()

			// когда middleware отработал, передаем данные далее в основную логику
			next.ServeHTTP(ww, r)
		}

		// Возвращаем созданный выше обработчик, приведя его к типу http.HandlerFunc
		return http.HandlerFunc(fn)
	}
}
