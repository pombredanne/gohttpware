package errors

import (
	"net/http"

	"github.com/tobi/airbrake-go"
	"github.com/zenazn/goji/web"
)

func AirbrakeRecoverer(apiKey string) func(*web.C, http.Handler) http.Handler {
	f := func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			airbrake.ApiKey = apiKey
			defer airbrake.CapturePanic(r)

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}
