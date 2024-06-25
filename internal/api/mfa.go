package api

import (
	"encoding/base32"
	"net/http"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/heyztb/lists/internal/crypto"
	"github.com/heyztb/lists/internal/database"
	"github.com/heyztb/lists/internal/log"
	"github.com/heyztb/lists/internal/middleware"
	"github.com/pquerna/otp/totp"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func VerifyTOTPCodeHandler(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading session context")
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := database.Users(
		database.UserWhere.ID.EQ(userID),
		qm.Load(database.UserRels.Setting),
	).One(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("error finding user from database")
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !user.R.Setting.MfaEnabled {
		render.Status(r, http.StatusBadRequest)
		w.Header().Add("HX-Redirect", "/login")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	b32NoPadding := base32.StdEncoding.WithPadding(base32.NoPadding)
	totpSecretBytes, err := crypto.AESDecrypt(crypto.ServerEncryptionKey, user.MfaSecret.Bytes)
	if err != nil {
		log.Err(err).Msgf("error decrypting totp secret for client %s", userID)
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	totpSecret := b32NoPadding.EncodeToString(totpSecretBytes)
	code := r.FormValue("code")
	if !totp.Validate(code, totpSecret) {
		log.Error().Msgf("bad totp code from client %s", userID)
		render.Status(r, http.StatusBadRequest)
		w.Header().Add("HX-Redirect", "/login")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	render.Status(r, http.StatusOK)
	w.Header().Add("HX-Redirect", "/app")
	w.WriteHeader(http.StatusOK)
}

func Disable2FAHandler(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading session context")
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := database.Users(
		database.UserWhere.ID.EQ(userID),
		qm.Load(database.UserRels.Setting),
	).One(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("error finding user from database")
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !user.R.Setting.MfaEnabled {
		render.Status(r, http.StatusBadRequest)
		w.Header().Add("HX-Redirect", "/app/settings")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.MfaSecret = null.BytesFrom(nil)
	user.MfaRecoveryCodes = []string{}
	user.R.Setting.MfaEnabled = false
	_, err = user.Update(r.Context(), database.DB, boil.Whitelist(
		database.UserColumns.MfaSecret,
		database.UserColumns.MfaRecoveryCodes,
	))
	if err != nil {
		log.Err(err).Msgf("error updating user in database for client %s", userID)
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = user.R.Setting.Update(r.Context(), database.DB, boil.Whitelist(
		database.SettingColumns.MfaEnabled,
	))
	if err != nil {
		log.Err(err).Msgf("error updating user settings in database for client %s", userID)
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("HX-Redirect", "/app/settings")
	w.WriteHeader(http.StatusNoContent)
}
