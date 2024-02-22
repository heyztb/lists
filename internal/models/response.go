package models

// response.go contains struct types for the responses returned by this API

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type IdentityResponse struct {
	Salt            string `json:"salt"`
	EphemeralPublic string `json:"B"`
}

type VerificationResponse struct {
	ServerProof string `json:"proof"`
}

type SuccessResponse struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}
