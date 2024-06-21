package html

import (
	"net/http"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/heyztb/lists/internal/database"
	"github.com/heyztb/lists/internal/html/templates/components/modals"
	"github.com/heyztb/lists/internal/html/templates/pages"
	"github.com/heyztb/lists/internal/html/templates/pages/app"
	"github.com/heyztb/lists/internal/log"
	"github.com/heyztb/lists/internal/middleware"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func ServeInternalServerErrorPage(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	pages.InternalServerError().Render(r.Context(), w)
}

func ServeNotFoundErrorPage(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotFound)
	pages.NotFoundErrorPage().Render(r.Context(), w)
}

func ServeMarketingIndex(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.Index("Lists").Render(r.Context(), w)
}

func ServeRegistration(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.Register().Render(r.Context(), w)
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	// we are adding this header here to redirect the client to this page
	// if they trigger the redirect in the auth middleware through an htmx request
	// i.e the button to show the change email/password modal(s) on the settings page
	// instead of loading the html for /login into the page it will redirect the browser to /login
	// this has no effect on visiting the page normally
	w.Header().Add("HX-Redirect", "/login")
	pages.Login().Render(r.Context(), w)
}

func ServePrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.PrivacyPolicy().Render(r.Context(), w)
}

func ServeTermsOfService(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.TermsOfService().Render(r.Context(), w)
}

func ServeAppIndex(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading context")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}

	user, err := database.Users(
		database.UserWhere.ID.EQ(userID),
		qm.Load(database.UserRels.Setting),
		qm.Load(database.UserRels.Lists),
		qm.Load(database.UserRels.Items),
		qm.Load(database.UserRels.Labels),
		qm.Load(database.UserRels.Comments),
	).One(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("error fetching user from database")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}

	render.Status(r, http.StatusOK)
	app.Index(*user).Render(r.Context(), w)
}

func ServeSettingsPage(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading context")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}

	user, err := database.FindUser(r.Context(), database.DB, userID)
	if err != nil {
		log.Err(err).Msg("error fetching user from database")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}
	render.Status(r, http.StatusOK)
	app.Settings(*user).Render(r.Context(), w)
}

func HTMXChangePasswordModal(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading context")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}
	log.Info().Msgf("popping change password modal for client %s", userID)
	render.Status(r, http.StatusOK)
	modals.ChangePassword().Render(r.Context(), w)
}

func HTMXChangeEmailModal(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading context")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}
	log.Info().Msgf("popping change email modal for client %s", userID)
	modals.ChangeEmail().Render(r.Context(), w)
}
