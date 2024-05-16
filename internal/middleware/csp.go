package middleware

import (
	"net/http"
)

// ContentSecurityPolicy middleware that provides the content security policy
// for this application. This is not the most stringent CSP possible, and
// certainly could use work to get it to a place where we are only executing
// code that I've written. This will do for now, however.
func ContentSecurityPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Security-Policy", "require-sri-for script; default-src 'none'; script-src 'self' 'unsafe-inline' 'unsafe-eval' 'wasm-unsafe-eval'; connect-src 'self'; img-src 'self'; style-src 'self' 'unsafe-hashes' 'sha256-9OfQztRBSnhT2ifc7/KgwOvhIpay6AeXqzSMt5gmEXk=' 'sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg='; frame-ancestors 'self'; form-action 'self';")
		next.ServeHTTP(w, r)
	})
}