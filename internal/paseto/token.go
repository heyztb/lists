package paseto

import (
	"fmt"
	"strconv"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/rs/zerolog/log"
)

var (
	issuer           = "lists-backend.api"
	audience         = "lists-frontend.client"
	ServerSigningKey paseto.V4AsymmetricSecretKey
)

func GenerateToken(userID uint64, expiration int, key []byte) (string, error) {
	now := time.Now()
	token := paseto.NewToken()
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(time.Duration(expiration) * time.Second))
	token.SetAudience(audience)
	token.SetIssuer(issuer)
	token.SetSubject(fmt.Sprint(userID))
	token.Set("dur", expiration)

	signedToken := token.V4Sign(ServerSigningKey, nil)

	return signedToken, nil
}

func ValidateToken(token string) (int64, int, error) {
	parser := paseto.MakeParser([]paseto.Rule{
		paseto.ValidAt(time.Now()),
		paseto.IssuedBy(issuer),
		paseto.ForAudience(audience),
	})
	parsedToken, err := parser.ParseV4Public(ServerSigningKey.Public(), token, nil)
	if err != nil {
		log.Err(err).Msg("failed to parse token")
		return -1, -1, fmt.Errorf("error parsing token: %w", err)
	}

	subject, err := parsedToken.GetSubject()
	if err != nil {
		log.Err(err).Msg("failed to get token subject")
		return -1, -1, fmt.Errorf("error getting token subject: %w", err)
	}

	userID, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		log.Err(err).Msg("failed to parse user ID")
		return -1, -1, fmt.Errorf("error parsing user ID: %w", err)
	}

	var sessionDuration int
	err = parsedToken.Get("dur", &sessionDuration)
	if err != nil {
		log.Err(err).Msg("failed to get session duration")
		return -1, -1, fmt.Errorf("failed to get session duration: %w", err)
	}

	return userID, sessionDuration, nil
}
