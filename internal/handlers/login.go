package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/1Password/srp"
	"github.com/go-chi/render"
	"github.com/heyztb/lists-backend/internal/cache"
	"github.com/heyztb/lists-backend/internal/database"
	"github.com/heyztb/lists-backend/internal/models"
	"github.com/heyztb/lists-backend/internal/paseto"
	"github.com/rs/zerolog/log"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("handler", "LoginHandler").Logger()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Err(err).Any("request", r).Msg("failed to read request body")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
	}
	req := &models.LoginRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		logger.Err(err).Bytes("body", body).Msg("failed to unmarshal body into login request struct")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}
	user, err := database.Users(
		database.UserWhere.Identifier.EQ(req.Identifier),
	).One(r.Context(), database.DB)
	if err != nil {
		logger.Err(err).Msg("failed to fetch user")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}
	srpServerBytes, err := cache.Cache.Get(
		fmt.Sprintf(cache.SRPServerKey, user.ID),
	)
	if err != nil {
		logger.Err(err).Msg("failed to fetch srp server from cache")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}
	srpServer := &srp.SRP{}
	if err = srpServer.UnmarshalBinary(srpServerBytes); err != nil {
		logger.Err(err).Msg("failed to unmarshal srp server bytes")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	s := srp.NumberFromString(user.Salt)
	proofBytes, err := hex.DecodeString(req.Proof)
	if err != nil {
		logger.Err(err).Msg("failed to decode client proof")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}
	if !srpServer.GoodServerProof(s.Bytes(), user.Identifier, proofBytes) {
		logger.Warn().Msg("failed to verify client proof")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}
	// we can ignore the error here because we will have already called this method before we get here (in the identity handler) therefore error will always be nil
	key, _ := srpServer.Key()
	expiration := user.R.Setting.SessionDuration
	err = cache.Redis.SetEx(
		r.Context(),
		fmt.Sprintf(cache.RedisSessionKeyPrefix, user.ID),
		hex.EncodeToString(key),
		time.Duration(expiration)*time.Second,
	).Err()
	if err != nil {
		logger.Err(err).Msg("failed to store session key in redis")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	token, err := paseto.GenerateToken(user.ID, expiration)
	if err != nil {
		logger.Err(err).Msg("failed to generate paseto token")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "lists-session",
		Value:    token,
		Path:     "/",
		Domain:   "localhost", // TODO: change this
		Expires:  time.Now().Add(time.Duration(expiration) * time.Second),
		Secure:   true,
		HttpOnly: true,
	})
	serverProof, _ := srpServer.M(s.Bytes(), user.Identifier)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &models.VerificationResponse{
		ServerProof: hex.EncodeToString(serverProof),
	})
}
