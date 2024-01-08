package header

import (
	"net/http"
	"time"

	"github.com/owncloud-ops/errors/pkg/version"
)

// Cache writes required cache headers to all requests.
func Cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		writer.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		writer.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

		next.ServeHTTP(writer, req)
	})
}

// Options writes required option headers to all requests.
func Options(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodOptions {
			next.ServeHTTP(writer, req)
		} else {
			writer.Header().Set("Access-Control-Allow-Origin", "*")
			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			writer.Header().Set("Allow", "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS")

			writer.WriteHeader(http.StatusOK)
		}
	})
}

// Secure writes required access headers to all requests.
func Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("X-Frame-Options", "DENY")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("X-XSS-Protection", "1; mode=block")

		if req.TLS != nil {
			writer.Header().Set("Strict-Transport-Security", "max-age=31536000")
		}

		next.ServeHTTP(writer, req)
	})
}

// Version writes the current API version to the headers.
func Version(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("X-ERRORS-VERSION", version.String)

		next.ServeHTTP(writer, req)
	})
}
