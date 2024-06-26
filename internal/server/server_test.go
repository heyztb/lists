// server_test contains the integration test suite for the backend
// this will likely be the single greatest test suite you've ever seen
// prepare yourself
package server_test

import (
	"bytes"
	"context"
	"crypto"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"code.posterity.life/srp/v2"
	"github.com/heyztb/lists/internal/log"
	"github.com/heyztb/lists/internal/models"
	"github.com/ory/dockertest/v3"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/argon2"
)

var (
	db          *sql.DB
	redisClient *redis.Client
	baseUrl     string
	httpClient  = http.Client{
		Transport: &http.Transport{},
	}
	salt      []byte
	srpClient *srp.Client
	triplet   srp.Triplet
)

func TestMain(m *testing.M) {
	log.Logger = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize pool")
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to ping docker")
	}

	redisContainer, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "backend-redis-test",
		Hostname:   "redis",
		Repository: "redis",
		Tag:        "latest",
		NetworkID:  "f03144698a89",
	})
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

	database, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "backend-db-test",
		Hostname:   "db",
		Repository: "backend-db",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=listsdb-testing",
			"POSTGRES_PASSWORD=testing",
			"POSTGRES_DB=lists-backend-test",
		},
		NetworkID: "f03144698a89",
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

	backend, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "backend-test",
		Repository: "listsbackend",
		Tag:        "latest",
		Env: []string{
			"LISTEN_ADDRESS=0.0.0.0:4322",
			"DISABLE_TLS=true",
			fmt.Sprintf("DATABASE_HOST=%s", "db"),
			fmt.Sprintf("DATABASE_PORT=%s", "5432"),
			"DATABASE_USER=listsdb-testing",
			"DATABASE_PASSWORD=testing",
			"DATABASE_NAME=lists-backend-test",
			"DATABASE_SSL_MODE=disable",
			fmt.Sprintf("REDIS_HOST=%s", "redis:6379"),
			"PASETO_KEY=5a6a2bd6c113a5087bf235b51474c1bf234e96c9417f3bb35417c698ceccaea3b527b1907e781650be31ad1108a9a12895e24331ca5687de1b6a7ee7e7363ad9",
		},
		Mounts: []string{
			"logs:/var/log/backend",
		},
		NetworkID: "f03144698a89",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create backend container")
	}

	if err := pool.Retry(func() error {
		addr := backend.GetHostPort("4322/tcp")
		baseUrl = fmt.Sprintf("http://%s", addr)
		_, err := http.Get(baseUrl)
		return err
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

func TestRegister(t *testing.T) {
	identifier := "hacker@hacker.com"
	password := "testing123"

	params := &srp.Params{
		Name:  "DH15-SHA256-Argon2",
		Group: srp.RFC5054Group3072,
		Hash:  crypto.SHA256,
		KDF: func(username string, password string, salt []byte) ([]byte, error) {
			p := []byte(username + ":" + password)
			key := argon2.IDKey(p, salt, 1, 64*1024, 4, 32)
			return key, nil
		},
	}

	var err error
	salt = srp.NewSalt()
	srpClient, err = srp.NewClient(params, identifier, password, salt)
	triplet, err = srp.ComputeVerifier(params, identifier, password, salt)
	assert.Nil(t, err, "failed to generate verifier")

	request := &models.RegistrationRequest{
		Identifier: identifier,
		Salt:       hex.EncodeToString(triplet.Salt()),
		Verifier:   hex.EncodeToString(triplet.Verifier()),
	}

	requestJson, err := json.Marshal(request)
	assert.Nil(t, err, "failed to marshal request json")

	var res *models.SuccessResponse
	err = makeRequest(http.MethodPost, "api/auth/register", bytes.NewReader(requestJson), &res)
	assert.Nil(t, err, "failed to make registration request")

	assert.Equal(t, http.StatusOK, res.Status)
	assert.Equal(t, "OK", res.Data)
}

func TestIdentity(t *testing.T) {
	A := srpClient.A()
	fmt.Println(hex.EncodeToString(A))

	request := &models.IdentityRequest{
		Identifier: "hacker@hacker.com",
	}

	requestJson, err := json.Marshal(request)
	assert.Nil(t, err, "failed to marshal request json")

	var res *models.IdentityResponse
	err = makeRequest(http.MethodPost, "api/auth/identify", bytes.NewReader(requestJson), &res)
	assert.Nil(t, err, "failed to make identity request")
	assert.Equal(t, http.StatusOK, res.Status)
	assert.NotEmpty(t, res.Salt)
	responseSalt, _ := hex.DecodeString(res.Salt)
	assert.Equal(t, salt, responseSalt)
	assert.NotEmpty(t, res.EphemeralPublic)
	B, err := hex.DecodeString(res.EphemeralPublic)
	assert.Nil(t, err, "failed to decode B")
	err = srpClient.SetB(B)
	assert.Nil(t, err, "invalid public key from server")
	fmt.Println("Salt: ", res.Salt, "B: ", res.EphemeralPublic)
}

func TestLogin(t *testing.T) {
	clientProof, err := srpClient.ComputeM1()
	assert.Nil(t, err, "failed to generate client proof")

	request := &models.LoginRequest{
		Identifier: "hacker@hacker.com",
		Proof:      hex.EncodeToString(clientProof),
	}

	requestJson, err := json.Marshal(request)
	assert.Nil(t, err, "failed to marshal request json")

	var res *models.LoginResponse
	err = makeRequest(http.MethodPost, "api/auth/login", bytes.NewReader(requestJson), &res)
	assert.Nil(t, err, "failed to make login request")
	assert.Equal(t, http.StatusOK, res.Status)
	assert.NotEmpty(t, res.ServerProof)

	serverProof, _ := hex.DecodeString(res.ServerProof)
	validProof, err := srpClient.CheckM2(serverProof)
	assert.Nil(t, err, "failed to check server proof")
	assert.True(t, validProof, "proof from server does not match what we expect")
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
			fmt.Println(string(responseBody))
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}
