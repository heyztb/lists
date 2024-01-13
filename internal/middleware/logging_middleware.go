package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			duration := time.Since(start)
			log.Info().
				Str("path", r.URL.Path).
				Dur("duration", duration).
				Msg("request info")
		}()

		next.ServeHTTP(w, r)
	})
}
