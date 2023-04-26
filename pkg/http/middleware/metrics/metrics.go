package metrics

import (
	"net/http"
	"time"
)

type metrics interface {
	IncrementTotalRequests()
	ObserveRequestDuration(t time.Duration)
}

func DurationMetrics(m metrics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(writer http.ResponseWriter, req *http.Request) {
			startedAt := time.Now()

			next.ServeHTTP(writer, req)

			m.IncrementTotalRequests()
			m.ObserveRequestDuration(time.Since(startedAt))
		}

		return http.HandlerFunc(fn)
	}
}
