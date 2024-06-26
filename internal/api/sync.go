package api

import (
	"encoding/base64"
	"net/http"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/heyztb/lists/internal/crypto"
	"github.com/heyztb/lists/internal/database"
	"github.com/heyztb/lists/internal/log"
	"github.com/heyztb/lists/internal/middleware"
	"github.com/heyztb/lists/internal/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func FullSyncHandler(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()

	userID, _, key, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

	queryMods := []qm.QueryMod{
		database.UserWhere.ID.EQ(userID),
		qm.Load(database.UserRels.Setting),
		qm.Load(database.UserRels.Lists, qm.Load(database.ListRels.Sections)),
		qm.Load(database.UserRels.Items),
		qm.Load(database.UserRels.Labels),
		qm.Load(database.UserRels.Comments),
	}

	user, err := database.Users(queryMods...).One(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("failed to load user from database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, user)
	if err != nil {
		log.Err(err).Msg("failed to encrypt user data")
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
		Data:   base64.RawStdEncoding.EncodeToString(encryptedJSON),
	})
}
