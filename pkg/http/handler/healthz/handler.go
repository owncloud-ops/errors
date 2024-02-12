package healthz

import (
	"io"
	"net/http"
)

// NewHandler creates handler for error pages serving.
func NewHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusOK)

		_, _ = io.WriteString(writer, http.StatusText(http.StatusOK))
	}
}
