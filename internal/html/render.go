package html

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"runtime"
	"time"

	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/heyztb/lists/internal/cache"
	"github.com/heyztb/lists/internal/crypto"
	"github.com/heyztb/lists/internal/database"
	"github.com/heyztb/lists/internal/html/templates/components/modals"
	"github.com/heyztb/lists/internal/html/templates/pages"
	"github.com/heyztb/lists/internal/html/templates/pages/app"
	"github.com/heyztb/lists/internal/log"
	"github.com/heyztb/lists/internal/middleware"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/argon2"
)

func ServeInternalServerErrorPage(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	pages.InternalServerError().Render(r.Context(), w)
}

func ServeNotFoundErrorPage(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotFound)
	pages.NotFoundErrorPage().Render(r.Context(), w)
}

func ServeMarketingIndex(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.Index("Lists").Render(r.Context(), w)
}

func ServeRegistration(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.Register().Render(r.Context(), w)
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	// we are adding this header here to redirect the client to this page
	// if they trigger the redirect in the auth middleware through an htmx request
	// i.e the button to show the change email/password modal(s) on the settings page
	// instead of loading the html for /login into the page it will redirect the browser to /login
	// this has no effect on visiting the page normally
	w.Header().Add("HX-Redirect", "/login")
	pages.Login().Render(r.Context(), w)
}

func ServePrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.PrivacyPolicy().Render(r.Context(), w)
}

func ServeTermsOfService(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	pages.TermsOfService().Render(r.Context(), w)
}

func ServeAppIndex(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading context")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}

	user, err := database.Users(
		database.UserWhere.ID.EQ(userID),
		qm.Load(database.UserRels.Setting),
		qm.Load(database.UserRels.Lists),
		qm.Load(database.UserRels.Items),
		qm.Load(database.UserRels.Labels),
		qm.Load(database.UserRels.Comments),
	).One(r.Context(), database.DB)
	if err != nil {
		log.Err(err).Msg("error fetching user from database")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}

	render.Status(r, http.StatusOK)
	app.Index(*user).Render(r.Context(), w)
}

func ServeSettingsPage(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading context")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}

	user, err := database.FindUser(r.Context(), database.DB, userID)
	if err != nil {
		log.Err(err).Msg("error fetching user from database")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}
	render.Status(r, http.StatusOK)
	app.Settings(*user).Render(r.Context(), w)
}

func HTMXChangePasswordModal(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading context")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}
	user, err := database.FindUser(r.Context(), database.DB, userID)
	if err != nil {
		log.Err(err).Msg("error fetching user from database")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}
	log.Info().Msgf("popping change password modal for client %s", userID)
	render.Status(r, http.StatusOK)
	modals.ChangePassword(user.Identifier).Render(r.Context(), w)
}

func HTMXChangeEmailModal(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(cmw.RequestIDKey).(string)
	log := log.Logger.With().Str("request_id", requestID).Logger()
	userID, _, _, err := middleware.ReadContext(r)
	if err != nil {
		log.Err(err).Msg("error reading context")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}
	user, err := database.FindUser(r.Context(), database.DB, userID)
	if err != nil {
		log.Err(err).Msg("error fetching user from database")
		render.Status(r, http.StatusInternalServerError)
		pages.InternalServerError().Render(r.Context(), w)
		return
	}
	log.Info().Msgf("popping change email modal for client %s", userID)
	modals.ChangeEmail(user.Identifier).Render(r.Context(), w)
}

func HTMXEnable2FAModal(w http.ResponseWriter, r *http.Request) {
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
	if user.R.Setting.MfaEnabled {
		render.Status(r, http.StatusBadRequest)
		w.Header().Add("HX-Redirect", "/app/settings")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Lists",
		AccountName: user.Identifier,
		Period:      30,
		SecretSize:  20,
		Digits:      6,
		Algorithm:   otp.AlgorithmSHA1,
		Rand:        rand.Reader,
	})
	if err != nil {
		log.Err(err).Msgf("error generating key for client %s", userID)
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		log.Err(err).Msgf("error generating key image for client %s", userID)
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = png.Encode(&buf, img)
	if err != nil {
		log.Err(err).Msgf("error encoding key image for client %s", userID)
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cache.Redis.SetEx(r.Context(), fmt.Sprintf("totp_secret:%s", userID), key.Secret(), time.Duration(1800)*time.Second)
	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	render.Status(r, http.StatusOK)
	modals.Enable2FA(key.Secret(), base64Image).Render(r.Context(), w)
}

func HTMX2FARecoveryCodesModal(w http.ResponseWriter, r *http.Request) {
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
	secret, err := cache.Redis.Get(r.Context(), fmt.Sprintf("totp_secret:%s", userID)).Result()
	if err != nil {
		log.Err(err).Msgf("error finding totp secret for client %s", userID)
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	code := r.FormValue("code")
	valid := totp.Validate(code, secret)
	if !valid {
		log.Warn().Msgf("invalid code received from client %s", userID)
		render.Status(r, http.StatusBadRequest)
		w.Header().Add("HX-Redirect", "/app/settings")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recoveryCodes := make([]string, 10)
	for i := range recoveryCodes {
		code, err := crypto.GenerateRandomString(12)
		if err != nil {
			log.Err(err).Msgf("error generating recovery code for client %s", userID)
			render.Status(r, http.StatusInternalServerError)
			w.Header().Add("HX-Redirect", "/500")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		recoveryCodes[i] = code
	}
	hashedCodes := make([]string, 10)
	for i := range recoveryCodes {
		salt := make([]byte, 12)
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			log.Err(err).Msgf("error generating salt for client %s", userID)
			render.Status(r, http.StatusInternalServerError)
			w.Header().Add("HX-Redirect", "/500")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		hashedCode := argon2.IDKey([]byte(recoveryCodes[i]), salt, 1, 64*1024, uint8(runtime.NumCPU()), 32)
		hashedCodes[i] = hex.EncodeToString(hashedCode)
	}
	b32NoPadding := base32.StdEncoding.WithPadding(base32.NoPadding)
	totpSecretBytes, err := b32NoPadding.DecodeString(secret)
	if err != nil {
		log.Err(err).Msg("error decoding base32 totp secret")
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encryptedTotpSecret, err := crypto.AESEncryptTableData(crypto.ServerEncryptionKey, totpSecretBytes)
	if err != nil {
		log.Err(err).Msgf("error encrypting totp secret for client %s", userID)
		render.Status(r, http.StatusInternalServerError)
		w.Header().Add("HX-Redirect", "/500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.MfaSecret.SetValid(encryptedTotpSecret)
	user.MfaRecoveryCodes = hashedCodes
	user.R.Setting.MfaEnabled = true
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
	cache.Redis.Del(r.Context(), fmt.Sprintf("totp_url:%s", userID))
	render.Status(r, http.StatusOK)
	modals.MFARecoveryCodes(recoveryCodes).Render(r.Context(), w)
}

func HTMXVerifyMFACode(w http.ResponseWriter, r *http.Request) {
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
}
