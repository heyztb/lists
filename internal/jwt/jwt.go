package jwt

import (
	"crypto/rsa"
	"fmt"
	"strconv"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

var (
	ServerSigningKey *rsa.PrivateKey
	issuer           = "lists.api"
	audience         = []string{"lists.client"}
)

type privateClaims struct {
	SessionDuration uint64 `json:"session_duration"`
}

func GenerateToken(userID uint64, expiration uint64) (string, error) {
	sig, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.RS256,
			Key:       ServerSigningKey,
		},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return "", err
	}

	t := jwt.Claims{
		Issuer:    issuer,
		Subject:   fmt.Sprintf("%d", userID),
		Audience:  audience,
		Expiry:    jwt.NewNumericDate(time.Now().Add(time.Duration(expiration))),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	privateClaims := privateClaims{
		SessionDuration: expiration,
	}

	token, err := jwt.Signed(sig).Claims(t).Claims(privateClaims).CompactSerialize()
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateToken(token string) (uint64, uint64, error) {
	t, err := jwt.ParseSigned(token)
	if err != nil {
		return 0, 0, err
	}

	registered := jwt.Claims{}
	private := privateClaims{}
	err = t.Claims(ServerSigningKey.PublicKey, &registered, &private)
	if err != nil {
		return 0, 0, err
	}

	err = registered.Validate(jwt.Expected{
		Issuer:   issuer,
		Audience: audience,
		Time:     time.Now(),
	})
	if err != nil {
		return 0, 0, err
	}

	userID, err := strconv.ParseInt(registered.Subject, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return uint64(userID), private.SessionDuration, nil
}
