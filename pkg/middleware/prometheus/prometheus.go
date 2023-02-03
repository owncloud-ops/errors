package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler initializes the prometheus middleware.
func Handler(token string) http.HandlerFunc {
	promHandler := promhttp.Handler()

	return func(writer http.ResponseWriter, req *http.Request) {
		if token == "" {
			promHandler.ServeHTTP(writer, req)

			return
		}

		header := req.Header.Get("Authorization")

		if header == "" {
			http.Error(writer, "Invalid or missing token", http.StatusUnauthorized)

			return
		}

		if header != "Bearer "+token {
			http.Error(writer, "Invalid or missing token", http.StatusUnauthorized)

			return
		}

		promHandler.ServeHTTP(writer, req)
	}
}
