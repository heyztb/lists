package server

import (
	"context"
	"crypto/tls"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-chi/chi/v5"
	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/heyztb/lists-backend/internal/database"
	"github.com/heyztb/lists-backend/internal/handlers"
	"github.com/heyztb/lists-backend/internal/middleware"
	security "github.com/heyztb/lists-backend/internal/paseto"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql/driver"
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
	dsn := driver.MySQLBuildQueryString(
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseSSLMode,
	)
	database.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	database.Redis = redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost,
		DB:   0,
	})

	if cfg.PasetoKey != "" {
		security.ServerSigningKey, err = paseto.NewV4AsymmetricSecretKeyFromHex(cfg.PasetoKey)
		if err != nil {
			log.Fatal().Err(err).Msg("faield to read paseto key")
		}
	}

	server := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: service(),
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
		serverStopCtx()
	}()

	if cfg.DisableTLS {
		err = server.ListenAndServe()
	} else {
		err = server.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile)
	}

	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("error starting server")
	}

	log.Info().Msgf("Server listening on %s", cfg.ListenAddress)

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func service() http.Handler {
	r := chi.NewRouter()
	r.Use(cmw.RequestID)
	r.Use(middleware.Logger)
	r.Use(cmw.Recoverer)
	r.Use(cmw.Heartbeat(`/`))
	r.Post(`/register`, handlers.EnrollmentHandler)
	r.Get(`/login`, handlers.IdentityHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Authentication)

		r.Get(`/lists`, handlers.GetListsHandler)
		r.Get(`/lists/{list}`, handlers.GetListHandler)
		r.Delete(`/lists/{list}`, handlers.DeleteListHandler)
		r.Get(`/sections`, handlers.GetSectionsHandler)
		r.Get(`/sections/{section}`, handlers.GetSectionHandler)
		r.Delete(`/sections/{section}`, handlers.DeleteSectionHander)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Decryption)
			r.Post(`/lists`, handlers.CreateListHandler)
			r.Post(`/lists/{list}`, handlers.UpdateListHandler)
			r.Post(`/sections`, handlers.CreateSectionHandler)
			r.Post(`/sections/{section}`, handlers.UpdateSectionHandler)
		})
	})

	return r
}
