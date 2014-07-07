package route

import (
	"net/http"
	"strings"
)

// Heartbeat endpoint middleware
func Heartbeat(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.EqualFold(r.URL.Path, "/heartbeat") {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("."))
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
