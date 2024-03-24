package paseto_test

import (
	"os"
	"testing"

	goPaseto "aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"github.com/heyztb/lists-backend/internal/paseto"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	paseto.ServerSigningKey = goPaseto.NewV4AsymmetricSecretKey()
	code := m.Run()
	os.Exit(code)
}

func TestGenerateToken(t *testing.T) {
	userID, err := uuid.NewV7()
	assert.Nil(t, err, "failed to generate uuid")
	token := paseto.GenerateToken(userID.String(), 3600)
	assert.NotEmpty(t, token, "failed to generate token")
}

func TestValidateToken(t *testing.T) {
	userUUID, err := uuid.NewV7()
	assert.Nil(t, err, "failed to generate uuid")
	token := paseto.GenerateToken(userUUID.String(), 3600)

	userID, duration, err := paseto.ValidateToken(token)
	assert.Nil(t, err, "failed to validate token")

	assert.Equal(t, userUUID.String(), userID, "user IDs do not match")
	assert.Equal(t, 3600, duration, "durations do not match")
}

func TestValidateTokenBadToken(t *testing.T) {
	// this token is signed with a different key than the one that was just generated above
	badToken := "v4.public.eyJhdWQiOiJsaXN0cy1mcm9udGVuZC5jbGllbnQiLCJkdXIiOjM2MDAsImV4cCI6IjIwMjQtMDMtMjRUMDA6NTY6MTUtMDc6MDAiLCJpYXQiOiIyMDI0LTAzLTIzVDIzOjU2OjE1LTA3OjAwIiwiaXNzIjoibGlzdHMtYmFja2VuZC5hcGkiLCJuYmYiOiIyMDI0LTAzLTIzVDIzOjU2OjE1LTA3OjAwIiwic3ViIjoiMDE4ZTZmM2UtZGY3Yy03NzBjLWEwYTUtODQxZTY0NTRkZTJjIn0QRaIyf3xaCQmFwipJo58PdGUQaDT985hFBgqKyj-1KJrR4FUB0MaQl8dmUmYmrkjqkrsCInx4iVqLTgjBUTsC"
	_, _, err := paseto.ValidateToken(badToken)
	assert.NotNil(t, err, "bad token passed as valid: something is wrong here")
}