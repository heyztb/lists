// Package middleware provides HTTP middleware functions for handling
// various aspects of HTTP request processing, such as logging, authentication,
// rate limiting, etc.
package middleware

import (
	"net/http"
	"time"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/heyztb/lists-backend/internal/log"
)

// Logger is an HTTP middleware that logs information about incoming requests.
// It logs the request ID, path, and duration.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Extract request ID from context
		requestID, ok := r.Context().Value(cmw.RequestIDKey).(string)
		if !ok {
			requestID = ""
		}

		defer func() {
			// Calculate request duration
			duration := time.Since(start)

			status, _ := r.Context().Value(render.StatusCtxKey).(int)

			// Log request information using structured logging
			log.Info().
				Str("request_id", requestID).
				Int("status", status).
				Str("path", r.URL.Path).
				Dur("duration", duration).
				Send()
		}()

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
