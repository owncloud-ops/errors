package metrics

import (
	"net/http"

	"github.com/owncloud-ops/errors/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewHandler(cfg *config.Config) http.HandlerFunc {
	promHandler := promhttp.HandlerFor(cfg.Metrics.Reg, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError})
	token := cfg.Metrics.Token

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
