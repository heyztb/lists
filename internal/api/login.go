package api

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"code.posterity.life/srp/v2"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/argon2"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/heyztb/lists/internal/cache"
	"github.com/heyztb/lists/internal/database"
	"github.com/heyztb/lists/internal/log"
	"github.com/heyztb/lists/internal/models"
	"github.com/heyztb/lists/internal/paseto"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
	req := &models.LoginRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		log.Err(err).Bytes("body", body).Msg("failed to unmarshal body into login request struct")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}
	user, err := database.Users(
		database.UserWhere.Identifier.EQ(req.Identifier),
		qm.Load(database.UserRels.Setting),
	).One(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("failed to fetch user")
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
		log.Err(err).Msg("failed to fetch srp server from cache")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}
	params := &srp.Params{
		Name:  "DH15-SHA256-Argon2",
		Group: srp.RFC5054Group3072,
		Hash:  crypto.SHA256,
		KDF: func(username string, password string, salt []byte) ([]byte, error) {
			p := []byte(username + ":" + password)
			key := argon2.IDKey(p, salt, 1, 64*1024, 4, 32)
			return key, nil
		},
	}
	srpServer, err := srp.RestoreServer(params, srpServerBytes)
	if err != nil {
		log.Err(err).Msg("failed to unmarshal srp server bytes")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	A, err := hex.DecodeString(req.EphemeralPublic)
	if err != nil {
		log.Err(err).Msg("failed to decode client ephemeral public key")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}
	err = srpServer.SetA(A)
	if err != nil {
		log.Err(err).Msg("bad client public key")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}
	proofBytes, err := hex.DecodeString(req.Proof)
	if err != nil {
		log.Err(err).Msg("failed to decode client proof")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}
	validProof, err := srpServer.CheckM1(proofBytes)
	if err != nil {
		log.Err(err).Msg("failed to check client proof")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	if !validProof {
		log.Warn().Msg("received invalid client proof")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}
	serverProof, err := srpServer.ComputeM2()
	if err != nil {
		log.Err(err).Msg("failed to generate server proof")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	key, err := srpServer.SessionKey()
	if err != nil {
		log.Err(err).Msg("failed to generate shared key")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	expiration := user.R.Setting.SessionDuration
	err = cache.Redis.SetEx(
		r.Context(),
		fmt.Sprintf(cache.RedisSessionKeyPrefix, user.ID),
		hex.EncodeToString(key),
		time.Duration(expiration)*time.Second,
	).Err()
	if err != nil {
		log.Err(err).Msg("failed to store session key in redis")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}
	token := paseto.GenerateToken(user.ID, expiration)
	http.SetCookie(w, &http.Cookie{
		Name:     "lists-session",
		Value:    token,
		Path:     "/",
		Domain:   "localhost", // TODO: change this
		Expires:  time.Now().Add(time.Duration(expiration) * time.Second),
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	})
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &models.LoginResponse{
		Status:      http.StatusOK,
		ServerProof: hex.EncodeToString(serverProof),
	})
}
