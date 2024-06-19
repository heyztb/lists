package middleware

import (
	"net/http"

	"github.com/go-chi/render"
)

// Dashboard middleware will redirect the user back to the app dashboard if
// they are logged in, otherwise we will allow the user to continue to their
// destination
func Dashboard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("lists-session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if err := sessionCookie.Valid(); err != nil {
			next.ServeHTTP(w, r)
			return
		}
		render.Status(r, http.StatusFound)
		http.Redirect(w, r, `/app`, http.StatusFound)
	})
}
