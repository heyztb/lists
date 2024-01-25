package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/heyztb/lists-backend/internal/database"
	"github.com/heyztb/lists-backend/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func CreateListHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &enrollmentResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to read request body",
		})
		return
	}

	log.Debug().Bytes("data", body).Msg("incoming request body")

	list := &models.List{}
	if err := json.Unmarshal(body, &list); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &enrollmentResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to unmarshal JSON body into List struct",
		})
		return
	}

	err = list.Insert(r.Context(), database.DB,
		boil.Whitelist(
			models.ListColumns.UserID,
			models.ListColumns.ParentID,
			models.ListColumns.Name,
			models.ListColumns.IsFavorite,
		),
	)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &enrollmentResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to save List to database",
		})
		return
	}
}
