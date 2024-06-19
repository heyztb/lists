package middleware

import "net/http"

// Size middleware will limit all incoming request body's to 32MB to prevent accidental or malicious requests from taking up resources
type Size struct {
	Mux http.Handler
}

func (s *Size) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20)
	s.Mux.ServeHTTP(w, r)
}
