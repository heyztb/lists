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

type enrollmentResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func EnrollmentHandler(w http.ResponseWriter, r *http.Request) {
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
	user := &models.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &enrollmentResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to unmarshal JSON body into User struct",
		})
		return
	}
	err = user.Insert(r.Context(), database.DB,
		boil.Whitelist(
			models.UserColumns.Username,
			models.UserColumns.Salt,
			models.UserColumns.Verifier,
		),
	)
	if err != nil {
		log.Error().Err(err).Msg("error inserting user")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &enrollmentResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to enroll user",
		})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &enrollmentResponse{
		Status:  http.StatusOK,
		Message: "Enrollment successful",
	})
}
