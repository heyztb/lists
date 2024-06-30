package html

import (
	"net/http"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/heyztb/lists/internal/database"
	"github.com/heyztb/lists/internal/log"
	"github.com/heyztb/lists/internal/middleware"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func FetchAvatarHandler(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading session context")
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := database.Users(
		database.UserWhere.ID.EQ(userID),
		qm.Load(database.UserRels.Setting),
	).One(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("error finding user from database")
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Cache-Control", "public, max-age=604800")
	render.Status(r, http.StatusOK)
	http.ServeFile(w, r, user.ProfilePicture.String)
}
