package html

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/render"
	"github.com/heyztb/lists-backend/internal/html/templates/pages"
)

// var ServeHomePage = templ.Handler(pages.Home("Lists")).ServeHTTP
var ServeRegisterPage = templ.Handler(pages.Register()).ServeHTTP
var ServeLoginPage = templ.Handler(pages.Login()).ServeHTTP

func ServeHomePage(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.Home("Lists").Render(r.Context(), w)
}