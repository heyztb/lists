package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/mail"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/heyztb/lists/internal/database"
	"github.com/heyztb/lists/internal/log"
	"github.com/heyztb/lists/internal/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Err(err).Any("request", r).Msg("failed to read request body")
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			render.Status(r, http.StatusRequestEntityTooLarge)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusRequestEntityTooLarge,
				Error:  "Content too large",
			})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	req := &models.RegistrationRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		log.Err(err).Bytes("body", body).Msg("failed to unmarshal body into registration request struct")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}
	email, err := mail.ParseAddress(req.Identifier)
	if err != nil {
		log.Err(err).Str("identifier", req.Identifier).Msg("error parsing identifier value as valid email address")
		render.Status(r, http.StatusBadRequest)
		w.Header().Add("HX-Redirect", "/app/settings")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := &database.User{
		Identifier: email.Address,
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
		log.Err(err).Msg("error inserting user")
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