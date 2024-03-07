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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("handler", "GetItemsHandler").Logger()

	userID, _, key, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

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

	filters := &models.GetItemsRequest{}
	if err := json.Unmarshal(body, &filters); err != nil {
		logger.Err(err).Bytes("body", body).Msg("failed to unmarshal request into GetItemsRequest struct")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	// by default we want to fetch the active items belonging to this user
	// if the request indicates that we want inactive items (i.e completed = true), then we fetch those instead
	queryMods := []qm.QueryMod{
		database.ItemWhere.CreatorID.EQ(userID),
		database.ItemWhere.IsCompleted.EQ(filters.Completed),
	}

	if filters.ListID != nil {
		queryMods = append(queryMods, database.ItemWhere.ListID.EQ(*filters.ListID))
	}

	if filters.SectionID != nil {
		queryMods = append(queryMods, database.ItemWhere.SectionID.EQ(null.Uint64From(*filters.SectionID)))
	}

	// TODO: implement this -- need to rethink how I handle labels
	// if filters.Label != "" {}

	items, err := database.Items(queryMods...).All(r.Context(), database.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusNotFound,
				Error:  "Not found",
			})
			return
		}
		logger.Err(err).Msg("error fetching items from database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, items)
	if err != nil {
		logger.Err(err).Msg("failed to encrypt section data")
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

func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("handler", "GetItemHandler").Logger()

	userID, _, key, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

	itemID := chi.URLParam(r, "item")
	itemIDInt, err := strconv.ParseInt(itemID, 10, 64)
	if err != nil {
		logger.Err(err).Str("item", itemID).Msg("invalid item ID")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	item, err := database.Items(
		database.ItemWhere.CreatorID.EQ(userID),
		database.ItemWhere.ID.EQ(uint64(itemIDInt)),
	).One(r.Context(), database.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusNotFound,
				Error:  "Not found",
			})
			return
		}
		logger.Err(err).Msg("error fetching item from database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, item)
	if err != nil {
		logger.Err(err).Msg("failed to encrypt section data")
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

func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("handler", "CreateItemHandler").Logger()

	userID, _, key, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

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

	request := &models.CreateItemRequest{}
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Err(err).Bytes("body", body).Msg("failed to unmarshal request into CreateItemRequest struct")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	// if no list_id param is sent in request then we default to the user's inbox project
	var listID uint64
	if request.ListID != nil {
		listID = *request.ListID
	} else {
		user, err := database.Users(
			qm.Load(database.UserRels.Lists),
			database.UserWhere.ID.EQ(userID),
		).One(r.Context(), database.DB)
		if err != nil {
			logger.Err(err).Msg("failed to load user and lists from database")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "Internal server error",
			})
			return
		}
		for _, list := range user.R.Lists {
			if list.IsInboxProject {
				listID = list.ID
			}
		}
	}

	item := &database.Item{
		ListID:      listID,
		SectionID:   null.Uint64FromPtr(request.SectionID),
		CreatorID:   userID,
		Content:     request.Content,
		Description: null.StringFromPtr(request.Description),
		ParentID:    null.Uint64FromPtr(request.ParentID),
		IsCompleted: false,
	}

	if request.Labels != nil {
		labelsJson, err := json.Marshal(*request.Labels)
		if err != nil {
			logger.Err(err).Msg("failed to marhsal labels json")
		} else {
			item.Labels = null.JSONFrom(labelsJson)
		}
	}

	if request.Priority != nil {
		item.Priority = *request.Priority
	} else {
		item.Priority = 1
	}

	if request.DueDate != nil {
		item.Due = null.TimeFromPtr(request.DueDate)
	} else if request.DueString != nil {
		// TODO: add support for parsing due_string values
	}

	if request.Duration != nil {
		item.Duration = null.IntFrom(int(*request.Duration))
	}

	if err := item.Insert(
		r.Context(),
		database.DB,
		boil.Blacklist(
			database.ItemColumns.ID,
			database.ItemColumns.CreatedAt,
			database.ItemColumns.UpdatedAt,
		),
	); err != nil {
		logger.Err(err).Msg("failed to insert item")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	if err = item.Reload(r.Context(), database.DB); err != nil {
		logger.Err(err).Msg("failed to reload item")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, item)
	if err != nil {
		logger.Err(err).Msg("failed to encrypt section data")
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

func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("handler", "UpdateItemHandler").Logger()

	userID, _, key, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

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

	itemID := chi.URLParam(r, "item")
	itemIDInt, err := strconv.ParseInt(itemID, 10, 64)
	if err != nil {
		logger.Err(err).Str("item", itemID).Msg("invalid item ID")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	request := &models.UpdateItemRequest{}
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Err(err).Bytes("body", body).Msg("failed to unmarshal request into UpdateItemRequest struct")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	item, err := database.Items(
		database.ItemWhere.ID.EQ(uint64(itemIDInt)),
		database.ItemWhere.CreatorID.EQ(userID),
	).One(r.Context(), database.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusNotFound,
				Error:  "Not found",
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

	err = request.UpdateItem(item)
	if err != nil {
		logger.Err(err).Msg("failed to update item")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	_, err = item.Update(r.Context(), database.DB, boil.Infer())
	if err != nil {
		logger.Err(err).Msg("failed to update item in database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, item)
	if err != nil {
		logger.Err(err).Msg("failed to encrypt item")
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

func CloseItemHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("handler", "CloseItemHandler").Logger()

	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

	itemID := chi.URLParam(r, "item")
	itemIDInt, err := strconv.ParseInt(itemID, 10, 64)
	if err != nil {
		logger.Err(err).Str("item", itemID).Msg("invalid item ID")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	item, err := database.Items(
		database.ItemWhere.ID.EQ(uint64(itemIDInt)),
		database.ItemWhere.CreatorID.EQ(userID),
	).One(r.Context(), database.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusNotFound,
				Error:  "Not found",
			})
			return
		}
		logger.Err(err).Msg("failed to get item")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	item.IsCompleted = true
	_, err = item.Update(r.Context(), database.DB, boil.Infer())
	if err != nil {
		logger.Err(err).Msg("failed to close item")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	render.NoContent(w, r)
}

func ReopenItemHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("handler", "ReopenItemHandler").Logger()

	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

	itemID := chi.URLParam(r, "item")
	itemIDInt, err := strconv.ParseInt(itemID, 10, 64)
	if err != nil {
		logger.Err(err).Str("item", itemID).Msg("invalid item ID")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	item, err := database.Items(
		database.ItemWhere.ID.EQ(uint64(itemIDInt)),
		database.ItemWhere.CreatorID.EQ(userID),
	).One(r.Context(), database.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusNotFound,
				Error:  "Not found",
			})
			return
		}
		logger.Err(err).Msg("failed to get item")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	item.IsCompleted = false
	_, err = item.Update(r.Context(), database.DB, boil.Infer())
	if err != nil {
		logger.Err(err).Msg("failed to reopen item")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	render.NoContent(w, r)
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Str("handler", "DeleteItemHandler").Logger()

	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

	itemID := chi.URLParam(r, "item")
	itemIDInt, err := strconv.ParseInt(itemID, 10, 64)
	if err != nil {
		logger.Err(err).Str("item", itemID).Msg("invalid item ID")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	item, err := database.Items(
		database.ItemWhere.ID.EQ(uint64(itemIDInt)),
		database.ItemWhere.CreatorID.EQ(userID),
	).One(r.Context(), database.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusNotFound,
				Error:  "Not found",
			})
			return
		}
		logger.Err(err).Msg("failed to get item")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	_, err = item.Delete(r.Context(), database.DB)
	if err != nil {
		logger.Err(err).Msg("failed to delete item")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	render.NoContent(w, r)
}
