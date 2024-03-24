// server_test contains the integration test suite for the backend
// this will likely be the single greatest test suite you've ever seen
// prepare yourself
package server_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/heyztb/lists-backend/internal/log"
	"github.com/heyztb/lists-backend/internal/models"
	"github.com/ory/dockertest/v3"
	"github.com/redis/go-redis/v9"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
)

var (
	db          *sql.DB
	redisClient *redis.Client
	baseUrl     string
	httpClient  = http.Client{
		Transport: &http.Transport{},
	}
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize pool")
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to ping docker")
	}

	redisContainer, err := pool.Run("redis", "latest", []string{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create redis container")
	}

	if err := pool.Retry(func() error {
		redisClient = redis.NewClient(&redis.Options{
			Addr: redisContainer.GetHostPort("6379/tcp"),
			DB:   0,
		})
		return redisClient.Ping(context.Background()).Err()
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}

	database, err := pool.Run("backend-db", "latest", []string{
		"POSTGRES_USER=listsdb-testing",
		"POSTGRES_PASSWORD=testing",
		"POSTGRES_DB=lists-backend-test",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database container")
	}

	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("user=listsdb-testing password=testing dbname=lists-backend-test host=127.0.0.1 port=%s sslmode=disable", database.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// run schema migrations
	migrations := migrate.FileMigrationSource{
		Dir: "../../sql",
	}

	migrationSlice, err := migrations.FindMigrations()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to find migrations")
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	if n != len(migrationSlice) {
		// wonder if we can ever reach this point, the docs seem to suggest so
		// but question is if the above err wil be nil or not if not all migrations run
		// too many levels of abstraction for me to want to look into it
		log.Fatal().Msg("did not run all migrations")
	}

	log.Info().Msgf("Applied %d migrations", n)

	backend, err := pool.Run("listsbackend", "latest", []string{
		"LISTEN_ADDRESS=0.0.0.0:4322",
		"DISABLE_TLS=true",
		fmt.Sprintf("DATABASE_HOST=%s", database.GetBoundIP("5432/tcp")),
		fmt.Sprintf("DATABASE_PORT=%s", database.GetPort("5432/tcp")),
		"DATABASE_USER=listsdb-testing",
		"DATABASE_PASSWORD=testing",
		"DATABASE_NAME=lists-backend-test",
		"DATABASE_SSL_MODE=disable",
		fmt.Sprintf("REDIS_HOST=%s", redisContainer.GetHostPort("6379/tcp")),
		"PASETO_KEY=5a6a2bd6c113a5087bf235b51474c1bf234e96c9417f3bb35417c698ceccaea3b527b1907e781650be31ad1108a9a12895e24331ca5687de1b6a7ee7e7363ad9",
	})

	if err := pool.Retry(func() error {
		addr := backend.GetHostPort("4322/tcp")
		baseUrl = fmt.Sprintf("http://%s", addr)
		_, err := http.Get(baseUrl)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to backend")
	}

	code := m.Run()

	if err := redisContainer.Close(); err != nil {
		log.Fatal().Err(err).Msg("could not close redis")
	}

	if err := database.Close(); err != nil {
		log.Fatal().Err(err).Msg("could not close database")
	}

	if err := backend.Close(); err != nil {
		log.Fatal().Err(err).Msg("could not close backend")
	}

	os.Exit(code)
}

// TestHealthcheck checks the healthcheck endpoint for the server
func TestHealthcheck(t *testing.T) {
	var res *models.SuccessResponse
	err := makeRequest(http.MethodGet, "", nil, &res)
	assert.Nil(t, err, "failed to make healthcheck request")
	assert.Equal(t, http.StatusOK, res.Status)
	assert.Equal(t, "OK", res.Data)
}

// makeRequest makes an http request to our server
func makeRequest(method, path string, body io.Reader, res interface{}) error {
	endpoint := fmt.Sprintf("%s/%s", baseUrl, path)
	url, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("bad url: %w", err)
	}

	request, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	defer response.Body.Close()

	if res != nil {
		err = json.Unmarshal(responseBody, res)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}