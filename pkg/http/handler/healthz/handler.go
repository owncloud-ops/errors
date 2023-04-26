package healthz

import (
	"io"
	"net/http"
)

// NewHandler creates handler for error pages serving.
func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
	}
}
