package static

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed all:assets
var assets embed.FS

func Mount(r chi.Router) {
	r.Route("/assets", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		})
		r.Handle("/*", http.FileServer(http.FS(assets)))
	})
}