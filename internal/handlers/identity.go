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
	"github.com/go-chi/render"
	"github.com/heyztb/lists-backend/internal/cache"
	"github.com/heyztb/lists-backend/internal/database"
	"github.com/heyztb/lists-backend/internal/models"
	"github.com/rs/zerolog/log"
)

func IdentityHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("handler", "IdentityHandler").Logger()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Err(err).Any("request", r).Msg("failed to read request body")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
	}
	req := &models.IdentityRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		logger.Err(err).Bytes("body", body).Msg("failed to unmarshal request into identity request struct")
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
			logger.Err(err).Msg("failed to fetch user from database")
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
		logger.Error().Msg("failed to initialize srp server component")
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
		logger.Err(err).Msg("invalid ephemeralPublicA from client")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	B := srpServer.EphemeralPublic()
	// eagerly generating the shared key now despite the user not being fully authenticated yet
	_, err = srpServer.Key()
	if err != nil {
		logger.Err(err).Msg("failed to generate shared key")
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
		logger.Err(err).Msg("failed to marshal srp server object to binary")
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
		logger.Err(err).Msg("failed to cache srp server object in memory")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &models.IdentityResponse{
		Salt:            user.Salt,
		EphemeralPublic: hex.EncodeToString(B.Bytes()),
	})
}
