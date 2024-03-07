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
	"github.com/heyztb/lists-backend/internal/log"
	"github.com/heyztb/lists-backend/internal/middleware"
	"github.com/heyztb/lists-backend/internal/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func GetSectionsHandler(w http.ResponseWriter, r *http.Request) {
	userID, _, key, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

	listID := r.URL.Query().Get("list_id")
	listIDInt, err := strconv.ParseInt(listID, 10, 64)
	if err != nil {
		log.Err(err).Str("list", listID).Msg("invalid list ID")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	sections, err := database.Sections(
		database.SectionWhere.UserID.EQ(userID),
		database.SectionWhere.ListID.EQ(uint64(listIDInt)),
	).All(r.Context(), database.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusNotFound,
				Error:  "Not found",
			})
			return
		}
		log.Err(err).Uint64("user_id", userID).Int64("list_id", listIDInt).Msg("failed to get sections from database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to get sections from database",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, sections)
	if err != nil {
		log.Err(err).Msg("failed to encrypt sections data")
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

func GetSectionHandler(w http.ResponseWriter, r *http.Request) {
	userID, _, key, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

	sectionID := chi.URLParam(r, "section")
	sectionIDInt, err := strconv.ParseInt(sectionID, 10, 64)
	if err != nil {
		log.Err(err).Str("section", sectionID).Msg("invalid section ID")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	section, err := database.Sections(
		database.SectionWhere.ID.EQ(uint64(sectionIDInt)),
		database.SectionWhere.UserID.EQ(userID),
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
		log.Err(err).Uint64("user_id", userID).Int64("section_id", sectionIDInt).Msg("failed to get section from database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, section)
	if err != nil {
		log.Err(err).Msg("failed to encrypt section data")
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

func CreateSectionHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Err(err).Any("request", r).Msg("failed to read request body")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	log.Debug().Bytes("data", body).Msg("incoming request body")

	request := &models.CreateSectionRequest{}
	if err := json.Unmarshal(body, &request); err != nil {
		log.Err(err).Bytes("body", body).Msg("failed to unmarshal request body into CreateSectionRequest struct")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	if request.Name == "" {
		log.Warn().Msg("blank name")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	section := &database.Section{
		UserID: userID,
		ListID: request.ListID,
		Name:   request.Name,
	}

	err = section.Insert(r.Context(), database.DB,
		boil.Whitelist(
			database.SectionColumns.UserID,
			database.SectionColumns.ListID,
			database.SectionColumns.Name,
		),
	)
	if err != nil {
		log.Err(err).Msg("failed to save section")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	err = section.Reload(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("failed to reload section")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, section)
	if err != nil {
		log.Err(err).Msg("failed to encrypt section data")
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

func UpdateSectionHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Err(err).Any("request", r).Msg("failed to read request body")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	sectionID := chi.URLParam(r, "section")
	sectionIDInt, err := strconv.ParseInt(sectionID, 10, 64)
	if err != nil {
		log.Err(err).Str("section", sectionID).Msg("invalid section ID")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	section, err := database.Sections(
		database.SectionWhere.ID.EQ(uint64(sectionIDInt)),
		database.SectionWhere.UserID.EQ(userID),
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
		log.Err(err).Uint64("user_id", userID).Str("section_id", sectionID).Msg("failed to fetch list from database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	request := &models.UpdateSectionRequest{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Err(err).Bytes("body", body).Uint64("user_id", userID).Str("section_id", sectionID).Msg("failed to unmarshal body into list")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	if request.Name == "" {
		log.Warn().Msg("blank name")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	section.Name = request.Name
	_, err = section.Update(r.Context(), database.DB, boil.Infer())
	if err != nil {
		log.Err(err).Msg("failed to update section in database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	err = section.Reload(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("failed to reload section from database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	encryptedJSON, err := crypto.AESEncrypt(key, section)
	if err != nil {
		log.Err(err).Msg("failed to encrypt section data")
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

func DeleteSectionHander(w http.ResponseWriter, r *http.Request) {
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized",
		})
		return
	}

	sectionID := chi.URLParam(r, "section")
	sectionIDInt, err := strconv.ParseInt(sectionID, 10, 64)
	if err != nil {
		log.Err(err).Str("section", sectionID).Msg("invalid section ID")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Bad request",
		})
		return
	}

	_, err = database.Sections(
		database.SectionWhere.ID.EQ(uint64(sectionIDInt)),
		database.SectionWhere.UserID.EQ(userID),
	).DeleteAll(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Int64("section_id", sectionIDInt).Msg("failed to delete list from database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
		return
	}

	render.NoContent(w, r)
}
