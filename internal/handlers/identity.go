package handlers

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/1Password/srp"
	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/heyztb/lists-backend/internal/cache"
	"github.com/heyztb/lists-backend/internal/database"
	"github.com/heyztb/lists-backend/internal/log"
	"github.com/heyztb/lists-backend/internal/models"
)

func IdentityHandler(w http.ResponseWriter, r *http.Request) {
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
	}
	req := &models.IdentityRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		log.Err(err).Bytes("body", body).Msg("failed to unmarshal request into identity request struct")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}
	user, err := database.Users(
		database.UserWhere.Identifier.EQ(req.Identifier),
	).One(r.Context(), database.DB)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Msg("failed to fetch user from database")
		}
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}
	v := srp.NumberFromString(user.Verifier)
	srpServer := srp.NewServerStd(srp.KnownGroups[srp.RFC5054Group3072], v)
	if srpServer == nil {
		log.Error().Msg("failed to initialize srp server component")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	A := srp.NumberFromString(req.EphemeralPublic)
	err = srpServer.SetOthersPublic(A)
	if err != nil {
		log.Err(err).Msg("invalid ephemeralPublicA from client")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}
	B := srpServer.EphemeralPublic()
	// eagerly generating the shared key now despite the user not being fully authenticated yet
	_, err = srpServer.Key()
	if err != nil {
		log.Err(err).Msg("failed to generate shared key")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	// marshal the srp server object into binary that way we are able to cache it in memory
	// for use later on -- this is important because we must maintain the same A and B values in order to generate and validate the key proof
	srpServerBytes, err := srpServer.MarshalBinary()
	if err != nil {
		log.Err(err).Msg("failed to marshal srp server object to binary")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	err = cache.Cache.Set(
		fmt.Sprintf(cache.SRPServerKey, user.ID),
		srpServerBytes,
	)
	if err != nil {
		log.Err(err).Msg("failed to cache srp server object in memory")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &models.IdentityResponse{
		Status:          http.StatusOK,
		Salt:            user.Salt,
		EphemeralPublic: hex.EncodeToString(B.Bytes()),
	})
}
