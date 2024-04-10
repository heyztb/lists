package html

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/heyztb/lists-backend/internal/html/templates/pages"
)

func ServeHomePage(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.Home("Lists").Render(r.Context(), w)
}

func ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.Register().Render(r.Context(), w)
}

func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.Login().Render(r.Context(), w)
}

func ServeAboutPage(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.About().Render(r.Context(), w)
}
