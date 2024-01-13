package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

func VerificationHandler(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "hello world")
}
