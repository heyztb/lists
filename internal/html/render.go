package html

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/heyztb/lists/internal/html/templates/pages"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.Home("Lists").Render(r.Context(), w)
}

func ServeRegistration(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.Register().Render(r.Context(), w)
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
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
