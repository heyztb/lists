package middleware

import (
	"net/http"
	"time"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestId, ok := r.Context().Value(cmw.RequestIDKey).(string)
		if !ok {
			requestId = ""
		}

		defer func() {
			duration := time.Since(start)
			log.Info().
				Str("request id", requestId).
				Str("path", r.URL.Path).
				Dur("duration", duration).
				Msg("request info")
		}()

		next.ServeHTTP(w, r)
	})
}
