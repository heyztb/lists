// package crypto provides high level functions for cryptographic routines used by this application
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/heyztb/lists/internal/log"
)

// generates random bytes to be used for a nonce value when encrypting messages to be sent to the client
func generateNonce() []byte {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Panic().Msg("failed to read random bytes for nonce")
	}

	return nonce
}

// AESEncrypt encrypts data using AES-256-GCM
func AESEncrypt(key []byte, data any) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("invalid key length: must be 32 bytes (256 bits) in length")
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("failed to marshal data to json for encryption")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("unable to create aes cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("unable to create gcm wrapped block cipher: %w", err)
	}

	nonce := generateNonce()
	sealed := gcm.Seal(nil, nonce, dataJSON, nil)

	return append(nonce, sealed...), nil
}

// AESDecrypt decrypts data using AES-256-GCM
func AESDecrypt(key []byte, data []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("invalid key length: must be 32 bytes (256 bits) in length")
	}

	nonce := data[:12]
	sealed := data[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("unable to create aes cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("unable to create gcm wrapped block cipher: %w", err)
	}

	unsealed, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt ciphertext: %w", err)
	}

	return unsealed, nil
}