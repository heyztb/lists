package main

import (
	"time"

	"github.com/JeremyLoy/config"
	"github.com/heyztb/lists-backend/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	// initialize the server config with some default values
	// these values can be overridden by environment variables
	cfg := &server.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	err := config.FromEnv().To(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read config from environment")
	}

	server.Run(cfg)
}