package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/heyztb/lists-backend/internal/models"
)

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &models.SuccessResponse{
		Status: http.StatusOK,
		Data:   "OK",
	})
}
