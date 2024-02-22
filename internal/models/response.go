package models

// response.go contains struct types for the responses returned by this API

type IdentityResponse struct {
	Salt            string `json:"salt"`
	EphemeralPublic string `json:"B"`
}

type VerificationResponse struct {
	ServerProof string `json:"proof"`
}
