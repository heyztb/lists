package server

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/hex"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-chi/chi/v5"
	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/heyztb/lists/internal/api"
	"github.com/heyztb/lists/internal/cache"
	"github.com/heyztb/lists/internal/crypto"
	"github.com/heyztb/lists/internal/database"
	"github.com/heyztb/lists/internal/html"
	"github.com/heyztb/lists/internal/html/static"
	"github.com/heyztb/lists/internal/log"
	"github.com/heyztb/lists/internal/middleware"
	security "github.com/heyztb/lists/internal/paseto"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql/driver"
)

type Config struct {
	// HTTP Server configuration
	ListenAddress string        `config:"LISTEN_ADDRESS"`
	ReadTimeout   time.Duration `config:"READ_TIMEOUT"`
	WriteTimeout  time.Duration `config:"WRITE_TIMEOUT"`
	IdleTimeout   time.Duration `config:"IDLE_TIMEOUT"`
	DisableTLS    bool          `config:"DISABLE_TLS"`
	TLSCertFile   string        `config:"TLS_CERT_FILE"`
	TLSKeyFile    string        `config:"TLS_KEY_FILE"`
	PasetoKey     string        `config:"PASETO_KEY"`
	AESKey        string        `config:"AES_KEY"`
	LogFilePath   string        `config:"LOG_FILE_PATH"`

	// Backing services configuration
	DatabaseHost     string `config:"DATABASE_HOST"`
	DatabasePort     int    `config:"DATABASE_PORT"`
	DatabaseUser     string `config:"DATABASE_USER"`
	DatabasePassword string `config:"DATABASE_PASSWORD"`
	DatabaseName     string `config:"DATABASE_NAME"`
	DatabaseSSLMode  string `config:"DATABASE_SSL_MODE"`
	RedisHost        string `config:"REDIS_HOST"`
}

func Run(cfg *Config) {
	var err error
	dsn := driver.PSQLBuildQueryString(
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseSSLMode,
	)
	database.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	cache.Redis = redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost,
		DB:   0,
	})

	if cfg.PasetoKey != "" {
		security.ServerSigningKey, err = paseto.NewV4AsymmetricSecretKeyFromHex(cfg.PasetoKey)
		if err != nil {
			log.Fatal().Err(err).Msg("faield to read paseto key")
		}
	}

	if cfg.AESKey != "" {
		crypto.ServerEncryptionKey, err = hex.DecodeString(cfg.AESKey)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to decode server AES key")
		}
	}

	server := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: &middleware.Size{Mux: service()},
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal().Msg("graceful shutdown timed out, forcing exit")
			}
			shutdownStopCtx()
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal().Err(err).Msg("error shutting down server")
		}
		log.Info().Msg("server shutting down")
		serverStopCtx()
	}()

	if cfg.DisableTLS {
		log.Info().Msgf("starting http server on %s", cfg.ListenAddress)
		server.ListenAndServe()
	} else {
		log.Info().Msgf("starting https server on %s", cfg.ListenAddress)
		server.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	os.Exit(0)
}

func service() http.Handler {
	r := chi.NewRouter()
	r.Use(cmw.RequestID)
	r.Use(middleware.Logger)
	r.Use(cmw.Recoverer)
	r.Use(middleware.ContentSecurityPolicy)
	r.NotFound(html.ServeNotFoundErrorPage)
	static.Mount(r)

	r.Get(`/privacy`, html.ServePrivacyPolicy)
	r.Get(`/tos`, html.ServeTermsOfService)
	r.Get(`/500`, html.ServeInternalServerErrorPage)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Dashboard)
		r.Get(`/`, html.ServeMarketingIndex)
		r.Get(`/register`, html.ServeRegistration)
		r.Get(`/login`, html.ServeLogin)
		r.Get(`/2fa`, html.HTMXVerifyMFACode)
	})

	r.Route(`/app`, func(r chi.Router) {
		r.Use(middleware.Authentication)
		r.Get(`/`, html.ServeAppIndex)
		r.Get(`/settings`, html.ServeSettingsPage)
	})

	r.Route(`/htmx`, func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authentication)
			r.Get(`/modal/changepassword`, html.HTMXChangePasswordModal)
			r.Get(`/modal/changeemail`, html.HTMXChangeEmailModal)
			r.Get(`/modal/enable2fa`, html.HTMXEnable2FAModal)
			r.Post(`/modal/2fa_recovery_codes`, html.HTMX2FARecoveryCodesModal)
		})
	})

	r.Route(`/api`, func(r chi.Router) {
		r.Get(`/`, api.HealthcheckHandler)
		r.Post(`/auth/register`, api.RegisterHandler)
		r.Post(`/auth/identify`, api.IdentityHandler)
		r.Post(`/auth/login`, api.LoginHandler)
		r.Post(`/auth/validate2fa`, api.VerifyTOTPCodeHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authentication)
			r.Post(`/auth/logout`, api.LogoutHandler)
			r.Delete(`/auth/delete`, api.DeleteAccountHandler)
			r.Delete(`/auth/disable2fa`, api.Disable2FAHandler)
			r.Patch(`/auth/updateverifier`, api.UpdateVerifierHandler)
			r.Patch(`/account/name`, api.UpdateNameHandler)
			r.Get(`/account/avatar`, html.FetchAvatarHandler)
			r.Patch(`/account/avatar`, api.UpdateAvatarHandler)

			r.Get(`/lists`, api.GetListsHandler)
			r.Get(`/lists/{list}`, api.GetListHandler)
			r.Delete(`/lists/{list}`, api.DeleteListHandler)
			r.Get(`/sections`, api.GetSectionsHandler)
			r.Get(`/sections/{section}`, api.GetSectionHandler)
			r.Delete(`/sections/{section}`, api.DeleteSectionHander)
			r.Get(`/items`, api.GetItemsHandler)
			r.Get(`/items/{item}`, api.GetItemHandler)
			r.Post(`/items/{item}/close`, api.CloseItemHandler)
			r.Post(`/items/{item}/reopen`, api.ReopenItemHandler)
			r.Delete(`/items/{item}`, api.DeleteItemHandler)
			r.Get(`/comments`, api.GetCommentsHandler)
			r.Get(`/comments/{comment}`, api.GetCommentHandler)
			r.Delete(`/comments/{comment}`, api.DeleteCommentHandler)
			r.Get(`/labels`, api.GetLabelsHandler)
			r.Get(`/labels/{label}`, api.GetLabelHandler)
			r.Delete(`/labels/{label}`, api.DeleteLabelHandler)

			r.Group(func(r chi.Router) {
				r.Use(middleware.Decryption)
				r.Post(`/lists`, api.CreateListHandler)
				r.Post(`/lists/{list}`, api.UpdateListHandler)
				r.Post(`/sections`, api.CreateSectionHandler)
				r.Post(`/sections/{section}`, api.UpdateSectionHandler)
				r.Post(`/items`, api.CreateItemHandler)
				r.Post(`/items/{item}`, api.UpdateItemHandler)
				r.Post(`/comments`, api.CreateCommentHandler)
				r.Post(`/comments/{comment}`, api.UpdateCommentHandler)
				r.Post(`/labels`, api.CreateLabelHandler)
				r.Post(`/labels/{label}`, api.UpdateLabelHandler)
			})
		})
	})

	return r
}
