package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/heyztb/lists-backend/internal/database"
	"github.com/heyztb/lists-backend/internal/jwt"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type authMiddlewareResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ctxKey struct {
	name string
}

var SessionKeyCtxKey = &ctxKey{"session-key"}

// Authentication middleware checks the incoming request for an Authorization header containing the requesting user's JWT. We validate this JWT by checking the signature as well as the issuer and audience.
// After parsing out the user ID and their configured session expiration time, we use that information
// to fetch and refresh the shared session key in redis. The key gets stored in the request context
// and then the request is passed onto the next handler in the chain.
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			userID, expiration, err := jwt.ValidateToken(token)
			if err != nil {
				log.Err(err).Msg("unable to validate jwt token")
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, &authMiddlewareResponse{
					Status:  http.StatusUnauthorized,
					Message: "Unauthorized",
				})
				return
			}
			storedKey, err := database.Redis.GetEx(
				r.Context(),
				fmt.Sprintf("%s:%d", database.RedisSessionKeyPrefix, userID),
				time.Duration(expiration)*time.Second,
			).Result()
			if err != nil {
				if err == redis.Nil {
					render.Status(r, http.StatusUnauthorized)
					render.JSON(w, r, &authMiddlewareResponse{
						Status:  http.StatusUnauthorized,
						Message: "Unauthorized",
					})
					return
				}
				log.Err(err).Msg("unable to fetch session key from redis")
			} else {
				r = r.WithContext(context.WithValue(r.Context(), SessionKeyCtxKey, storedKey))
			}
		}
		next.ServeHTTP(w, r)
	})
}
