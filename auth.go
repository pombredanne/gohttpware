package httpware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func Unauthorized(realm string, w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Please authenticate with the proper API user and password!\n"))
}

func BasicAuth(user, pass, realm string) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		authString := fmt.Sprintf("%s:%s", user, pass)
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Basic ") {
				Unauthorized(realm, w)
				return
			}

			pass, err := decodeAuth(auth[6:])
			if err != nil || pass != authString {
				Unauthorized(realm, w)
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}

// This function decodes the given string.
// Here is where we would put any decryption if required.
func decodeAuth(auth string) (string, error) {
	pass, err := base64.StdEncoding.DecodeString(auth)
	return string(pass), err
}
