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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("handler", "RegisterHandler").Logger()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Err(err).Any("request", r).Msg("failed to read request body")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	req := &models.RegistrationRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		logger.Err(err).Bytes("body", body).Msg("failed to unmarshal body into registration request struct")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}
	user := &database.User{
		Identifier: req.Identifier,
		Salt:       req.Salt,
		Verifier:   req.Verifier,
	}
	err = user.Insert(r.Context(), database.DB,
		boil.Whitelist(
			database.UserColumns.Identifier,
			database.UserColumns.Salt,
			database.UserColumns.Verifier,
		),
	)
	if err != nil {
		logger.Err(err).Msg("error inserting user")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &models.SuccessResponse{
		Status: http.StatusOK,
		Data:   "OK",
	})
}