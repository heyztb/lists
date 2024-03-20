package models

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type IdentityResponse struct {
	Status          int    `json:"status"`
	Salt            string `json:"salt"`
	EphemeralPublic string `json:"B"`
}

type LoginResponse struct {
	Status      int    `json:"status"`
	ServerProof string `json:"proof"`
}

type SuccessResponse struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}
