// Package middleware provides HTTP middleware functions for handling
// various aspects of HTTP request processing, such as logging, authentication,
// rate limiting, etc.
package middleware

import (
	"net/http"
	"time"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/heyztb/lists-backend/internal/log"
)

// Logger is an HTTP middleware that logs information about incoming requests.
// It logs the request ID, request path, and request duration.
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

			// Log request information using structured logging
			log.Info().
				Str("id", requestID).
				Str("path", r.URL.Path).
				Dur("duration", duration).
				Msg("request info")
		}()

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
