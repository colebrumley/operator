package api

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// BasicAuth accepts the user:pass string and a handler function. If the user
// is successfully validated, the request is passed through to the provided handler.
func BasicAuth(pw string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header.Get("Authorization")) <= 0 {
			http.Error(w, "authentication is required", http.StatusUnauthorized)
			return
		}
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
		if auth[0] != "Basic" || len(auth) != 2 {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		parsed := string(payload)
		if strings.Contains(parsed, ":") {
			pair := strings.SplitN(string(payload), ":", 2)
			parsed = pair[1]
		}
		if !Validate(pw, parsed) {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

// Validate does the comparison between the provided password and the valid password
func Validate(pw string, test string) bool {
	if test == pw {
		return true
	}
	return false
}
