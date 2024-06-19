package api

import (
	"fmt"
	"net/http"
	"time"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/heyztb/lists/internal/cache"
	"github.com/heyztb/lists/internal/log"
	"github.com/heyztb/lists/internal/middleware"
	"github.com/heyztb/lists/internal/models"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading session context on logout")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}
	err = cache.Redis.Del(
		r.Context(),
		fmt.Sprintf(cache.RedisSessionKeyPrefix, userID),
	).Err()
	if err != nil {
		log.Err(err).Msg("error deleting shared key from redis")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "lists-session",
		Value:    "",
		Path:     "/",
		Domain:   "localhost", // TODO: change this
		Expires:  time.Unix(0, 0),
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	})
	render.Status(r, http.StatusFound)
	w.WriteHeader(http.StatusFound)
}
