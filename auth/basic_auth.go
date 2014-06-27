package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

type AuthFunc func(*http.Request, []string) bool

func Wrap(f AuthFunc, realm string, secrets ...string) func(http.Handler) http.Handler {
	h := func(h http.Handler) http.Handler {
		hn := func(w http.ResponseWriter, r *http.Request) {
			if !f(r, secrets) {
				Unauthorized(realm, w)
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hn)
	}
	return h
}

func Unauthorized(realm string, w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Please authenticate with the proper details\n"))
}

func BasicAuth(user, pass string) func(http.Handler) http.Handler {
	return Wrap(basicAuth, "Restricted", user, pass)
}

func basicAuth(r *http.Request, secrets []string) bool {
	user, pass := secrets[0], secrets[1]
	authString := fmt.Sprintf("%s:%s", user, pass)
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		return false
	}

	pass, err := decodeAuth(auth[6:])
	if err != nil || pass != authString {
		return false
	}
	return true
}

// This function decodes the given string.
// Here is where we would put any decryption if required.
func decodeAuth(auth string) (string, error) {
	pass, err := base64.StdEncoding.DecodeString(auth)
	return string(pass), err
}
