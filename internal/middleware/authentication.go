package middleware

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/heyztb/lists-backend/internal/cache"
	"github.com/heyztb/lists-backend/internal/log"
	"github.com/heyztb/lists-backend/internal/models"
	"github.com/heyztb/lists-backend/internal/paseto"
)

type ctxKey struct {
	name string
}

var (
	UserIDCtxKey          = &ctxKey{"user-id"}
	SessionDurationCtxKey = &ctxKey{"session-duration"}
	SessionKeyCtxKey      = &ctxKey{"session-key"}
)

// Authentication middleware checks the incoming request for the presence of a session cookie as set by the VerificationHandler
// Upon confirmation of the cookie, we validate and parse the token contained within so that we can populate the request context
// All subsequent handlers will have access to the user ID, session duration, and shared session key.
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
		log := log.Logger.With().Str("request_id", requestID).Logger()

		sessionCookie, err := r.Cookie("lists-session")
		if err != nil {
			log.Err(err).Msg("unable to get session cookie")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusUnauthorized,
				Error:  "Unauthorized",
			})
			return
		}
		if err := sessionCookie.Valid(); err != nil {
			log.Err(err).Msg("unable to validate session cookie")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusUnauthorized,
				Error:  "Unauthorized",
			})
			return
		}
		userID, expiration, err := paseto.ValidateToken(sessionCookie.Value)
		if err != nil {
			log.Err(err).Msg("unable to validate jwt token")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusUnauthorized,
				Error:  "Unauthorized",
			})
			return
		}
		storedKey, err := cache.Redis.GetEx(
			r.Context(),
			fmt.Sprintf(cache.RedisSessionKeyPrefix, userID),
			time.Duration(expiration)*time.Second,
		).Result()
		if err != nil {
			log.Err(err).Msg("unable to fetch session key from redis")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusUnauthorized,
				Error:  "Unauthorized",
			})
			return
		}
		r = populateContext(r, map[*ctxKey]any{
			UserIDCtxKey:          userID,
			SessionDurationCtxKey: expiration,
			SessionKeyCtxKey:      storedKey,
		})
		next.ServeHTTP(w, r)
	})
}

func populateContext(r *http.Request, values map[*ctxKey]any) *http.Request {
	for k, v := range values {
		// store raw session key in context to save handlers from having to do it
		// ignore the error because it will never happen
		if k == SessionKeyCtxKey {
			v, _ = hex.DecodeString(v.(string))
		}
		r = r.WithContext(context.WithValue(r.Context(), k, v))
	}
	return r
}

func ReadContext(r *http.Request) (uint64, uint64, []byte, error) {
	userID, ok := r.Context().Value(UserIDCtxKey).(uint64)
	if !ok {
		return 0, 0, nil, errors.New("no user ID in request context")
	}
	expiration, ok := r.Context().Value(SessionDurationCtxKey).(uint64)
	if !ok {
		return 0, 0, nil, errors.New("no session duration in request context")
	}
	storedKey, ok := r.Context().Value(SessionKeyCtxKey).([]byte)
	if !ok {
		return 0, 0, nil, errors.New("no session key in request context")
	}

	return userID, expiration, storedKey, nil
}
