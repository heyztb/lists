package paseto

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

var (
	issuer           = "lists-backend.api"
	audience         = "lists-frontend.client"
	ServerSigningKey paseto.V4AsymmetricSecretKey
)

func GenerateToken(userID string, expiration int) string {
	now := time.Now()
	token := paseto.NewToken()
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(time.Duration(expiration) * time.Second))
	token.SetAudience(audience)
	token.SetIssuer(issuer)
	token.SetSubject(userID)
	token.Set("dur", expiration)

	signedToken := token.V4Sign(ServerSigningKey, nil)

	return signedToken
}

func ValidateToken(token string) (string, int, error) {
	parser := paseto.MakeParser([]paseto.Rule{
		paseto.ValidAt(time.Now()),
		paseto.IssuedBy(issuer),
		paseto.ForAudience(audience),
	})
	parsedToken, err := parser.ParseV4Public(ServerSigningKey.Public(), token, nil)
	if err != nil {
		return "", -1, fmt.Errorf("error parsing token: %w", err)
	}

	subject, err := parsedToken.GetSubject()
	if err != nil {
		return "", -1, fmt.Errorf("error getting token subject: %w", err)
	}

	var sessionDuration int
	err = parsedToken.Get("dur", &sessionDuration)
	if err != nil {
		return "", -1, fmt.Errorf("failed to get session duration: %w", err)
	}

	return subject, sessionDuration, nil
}
