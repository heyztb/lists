package crypto_test

import (
	"encoding/base64"
	"encoding/hex"
	"os"
	"testing"

	"github.com/heyztb/lists/internal/crypto"
	"github.com/stretchr/testify/assert"
)

var key = make([]byte, 32)

func TestMain(m *testing.M) {
	key, _ = hex.DecodeString("88c46c6c9944ea996157b96a0c67115d2f7ec2a9818d9e92d69bf9640689ff01")
	code := m.Run()
	os.Exit(code)
}

func TestAESEncrypt(t *testing.T) {
	data := map[string]string{"hello": "world"}
	encrypted, err := crypto.AESEncrypt(key, data)
	assert.Nil(t, err, "failed to encrypt data")
	assert.NotEmpty(t, encrypted)
}

func TestAESDecrypt(t *testing.T) {
	encrypted := "WjmC28+UlwFru2gDCtMjGuVkZ5MZ4av7H4RHQCWtFZbOLrvHC9Y2sYo64PCh"
	decoded, _ := base64.RawStdEncoding.DecodeString(encrypted)
	decrypted, err := crypto.AESDecrypt(key, decoded)
	assert.Nil(t, err, "failed to decrypt")
	assert.NotEmpty(t, decrypted)
	assert.Equal(t, `{"hello":"world"}`, string(decrypted))
}
