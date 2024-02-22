package models

// requests.go contains struct types for incoming request bodies

type EnrollmentRequest struct {
	Identifier string `json:"identifier"`
	Salt       string `json:"salt"` // hex
	Verifier   string `json:"verifier"`
}

type IdentityRequest struct {
	Identifier      string `json:"identifier"`
	EphemeralPublic string `json:"A"`
}

type VerificationRequest struct {
	Identifier string `json:"identifier"`
	Proof      string `json:"proof"`
}
