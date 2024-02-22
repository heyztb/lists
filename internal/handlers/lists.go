package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/heyztb/lists-backend/internal/crypto"
	"github.com/heyztb/lists-backend/internal/database"
	"github.com/heyztb/lists-backend/internal/middleware"
	"github.com/heyztb/lists-backend/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func CreateListHandler(w http.ResponseWriter, r *http.Request) {
	userID, _, key, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to read request body",
		})
		return
	}

	log.Debug().Bytes("data", body).Msg("incoming request body")

	list := &database.List{}
	if err := json.Unmarshal(body, &list); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to unmarshal JSON body into List struct",
		})
		return
	}
	list.UserID = userID

	err = list.Insert(r.Context(), database.DB,
		boil.Whitelist(
			database.ListColumns.UserID,
			database.ListColumns.ParentID,
			database.ListColumns.Name,
			database.ListColumns.IsFavorite,
		),
	)
	if err != nil {
		log.Err(err).Msg("failed to save list")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to save List to database",
		})
		return
	}

	err = list.Reload(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("failed to reload list")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to reload List from database",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, list)
	if err != nil {
		log.Err(err).Msg("failed to encrypt list data")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &models.SuccessResponse{
		Status: http.StatusOK,
		Data:   base64.RawStdEncoding.EncodeToString(encryptedJSON),
	})
}

func UpdateListHandler(w http.ResponseWriter, r *http.Request) {
	userID, _, key, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to read request body",
		})
		return
	}

	listID := chi.URLParam(r, "list")
	listIDInt, err := strconv.ParseInt(listID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid list ID",
		})
		return
	}

	list, err := database.Lists(
		database.ListWhere.ID.EQ(uint64(listIDInt)),
		database.ListWhere.UserID.EQ(userID),
	).One(r.Context(), database.DB)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Uint64("user_id", userID).Str("list_id", listID).Msg("failed to fetch list from database")
		}
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to fetch list from database",
		})
		return
	}

	err = json.Unmarshal(body, &list)
	if err != nil {
		log.Err(err).Bytes("body", body).Uint64("user_id", userID).Str("list_id", listID).Msg("failed to unmarshal body into list")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to parse json into list",
		})
		return
	}

	_, err = list.Update(r.Context(), database.DB, boil.Whitelist(
		database.ListColumns.Name,
		database.ListColumns.IsFavorite,
	))
	if err != nil {
		log.Err(err).Msg("failed to update list in database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to update list",
		})
		return
	}

	err = list.Reload(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("failed to reload list from database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to reload list",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, list)
	if err != nil {
		log.Err(err).Msg("failed to encrypt list data")
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

func DeleteListHandler(w http.ResponseWriter, r *http.Request) {
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

	listID := chi.URLParam(r, "list")
	listIDInt, err := strconv.ParseInt(listID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid list ID",
		})
		return
	}

	_, err = database.Lists(
		database.ListWhere.ID.EQ(uint64(listIDInt)),
		database.ListWhere.UserID.EQ(userID),
	).DeleteAll(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Int64("list_id", listIDInt).Msg("failed to delete list from database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to delete list",
		})
		return
	}

	render.Status(r, http.StatusNoContent)
	render.JSON(w, r, struct{}{})
}

func GetListsHandler(w http.ResponseWriter, r *http.Request) {
	userID, _, key, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

	lists, err := database.Lists(
		database.ListWhere.UserID.EQ(userID),
	).All(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("failed to fetch list from database")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to fetch lists from database",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, lists)
	if err != nil {
		log.Err(err).Msg("failed to encrypt lists data")
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

func GetListHandler(w http.ResponseWriter, r *http.Request) {
	userID, _, key, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

	listID := chi.URLParam(r, "list")
	listIDInt, err := strconv.ParseInt(listID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid list ID",
		})
		return
	}

	list, err := database.Lists(
		database.ListWhere.ID.EQ(uint64(listIDInt)),
		database.ListWhere.UserID.EQ(userID),
	).One(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("failed to fetch list from database")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to fetch list from database",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, list)
	if err != nil {
		log.Err(err).Msg("failed to encrypt list data")
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
