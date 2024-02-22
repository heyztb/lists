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
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid list ID",
		})
		return
	}

	sections, err := database.Sections(
		database.SectionWhere.UserID.EQ(userID),
		database.SectionWhere.ListID.EQ(uint64(listIDInt)),
	).All(r.Context(), database.DB)
	if err != nil {
		if err != sql.ErrNoRows {
			render.Status(r, http.StatusNoContent)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusNoContent,
				Error:  "No sections found",
			})
			return
		} else {
			log.Err(err).Uint64("user_id", userID).Int64("list_id", listIDInt).Msg("failed to get sections from database")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "Failed to get sections from database",
			})
			return
		}
	}

	encryptedJSON, err := crypto.AESEncrypt(key, sections)
	if err != nil {
		log.Err(err).Msg("failed to encrypt sections data")
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
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid section ID",
		})
		return
	}

	section, err := database.Sections(
		database.SectionWhere.ID.EQ(uint64(sectionIDInt)),
		database.SectionWhere.UserID.EQ(userID),
	).One(r.Context(), database.DB)
	if err != nil {
		if err != sql.ErrNoRows {
			render.Status(r, http.StatusNoContent)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusNoContent,
				Error:  "No section found",
			})
			return
		} else {
			log.Err(err).Uint64("user_id", userID).Int64("section_id", sectionIDInt).Msg("failed to get section from database")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, &models.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "Failed to get section from database",
			})
			return
		}
	}

	encryptedJSON, err := crypto.AESEncrypt(key, section)
	if err != nil {
		log.Err(err).Msg("failed to encrypt section data")
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
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to read request body",
		})
		return
	}

	log.Debug().Bytes("data", body).Msg("incoming request body")

	section := &database.Section{}
	if err := json.Unmarshal(body, &section); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to unmarshal JSON body into Section struct",
		})
		return
	}
	section.UserID = userID

	err = section.Insert(r.Context(), database.DB,
		boil.Whitelist(
			database.SectionColumns.UserID,
			database.SectionColumns.ListID,
			database.SectionColumns.Name,
			database.SectionColumns.Position,
		),
	)
	if err != nil {
		log.Err(err).Msg("failed to save section")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to save section to database",
		})
		return
	}

	err = section.Reload(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("failed to reload section")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to reload section from database",
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
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to read request body",
		})
		return
	}

	sectionID := chi.URLParam(r, "section")
	sectionIDInt, err := strconv.ParseInt(sectionID, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid section ID",
		})
		return
	}

	section, err := database.Sections(
		database.SectionWhere.ID.EQ(uint64(sectionIDInt)),
		database.SectionWhere.UserID.EQ(userID),
	).One(r.Context(), database.DB)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Uint64("user_id", userID).Str("section_id", sectionID).Msg("failed to fetch list from database")
		}
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to fetch list from database",
		})
		return
	}

	err = json.Unmarshal(body, &section)
	if err != nil {
		log.Err(err).Bytes("body", body).Uint64("user_id", userID).Str("section_id", sectionID).Msg("failed to unmarshal body into list")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Failed to parse json into list",
		})
		return
	}

	_, err = section.Update(r.Context(), database.DB, boil.Whitelist(
		database.SectionColumns.Name,
	))
	if err != nil {
		log.Err(err).Msg("failed to update section in database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to update section",
		})
		return
	}

	err = section.Reload(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("failed to reload section from database")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to reload section",
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
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid section ID",
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
			Error:  "Failed to delete list",
		})
		return
	}

	render.Status(r, http.StatusNoContent)
	render.JSON(w, r, struct{}{})
}
