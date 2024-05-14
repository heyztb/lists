package models

type HealthcheckResponse struct {
	Status   int    `json:"status"`
	Checksum string `json:"checksum"`
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

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
