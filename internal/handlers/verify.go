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
	"github.com/heyztb/lists-backend/internal/database"
	"github.com/heyztb/lists-backend/internal/models"
	"github.com/heyztb/lists-backend/internal/paseto"
	"github.com/rs/zerolog/log"
)

func VerificationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &response{
			Status:  http.StatusBadRequest,
			Message: "Bad request",
		})
	}
	req := &models.VerificationRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &response{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}
	user, err := database.Users(
		database.UserWhere.Identifier.EQ(req.Identifier),
	).One(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("failed to fetch user")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &response{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}
	srpServerBytes, err := database.Cache.Get(
		fmt.Sprintf(database.SRPServerKey, user.ID),
	)
	if err != nil {
		log.Err(err).Msg("failed to fetch srp server from cache")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &response{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}
	srpServer := &srp.SRP{}
	if err = srpServer.UnmarshalBinary(srpServerBytes); err != nil {
		log.Err(err).Msg("failed to unmarshal srp server bytes")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &response{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}
	s := srp.NumberFromString(user.Salt)
	proofBytes, err := hex.DecodeString(req.Proof)
	if err != nil {
		log.Err(err).Msg("failed to decode client proof")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &response{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}
	if !srpServer.GoodServerProof(s.Bytes(), user.Identifier, proofBytes) {
		log.Warn().Msg("failed to verify client proof")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &response{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}
	// we can ignore the error here because we will have already called this method before we get here (in the identity handler) therefore error will always be nil
	key, _ := srpServer.Key()
	expiration := user.R.Setting.SessionDuration
	err = database.Redis.SetEx(
		r.Context(),
		fmt.Sprintf(database.RedisSessionKeyPrefix, user.ID),
		hex.EncodeToString(key),
		time.Duration(expiration)*time.Second,
	).Err()
	if err != nil {
		log.Err(err).Msg("failed to store session key in redis")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &response{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}
	token, err := paseto.GenerateToken(user.ID, expiration, key)
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
