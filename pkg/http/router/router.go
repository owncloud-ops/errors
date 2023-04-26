package router

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.owncloud.com/owncloud-ops/errors/pkg/config"
	errorpagesHandler "github.owncloud.com/owncloud-ops/errors/pkg/http/handler/errorpage"
	healthHandler "github.owncloud.com/owncloud-ops/errors/pkg/http/handler/healthz"
	metricsHandler "github.owncloud.com/owncloud-ops/errors/pkg/http/handler/metrics"
	"github.owncloud.com/owncloud-ops/errors/pkg/http/handler/notfound"
	"github.owncloud.com/owncloud-ops/errors/pkg/http/middleware/header"
	"github.owncloud.com/owncloud-ops/errors/pkg/http/middleware/metrics"
)

const MiddlewareTimeout = 60 * time.Second

// Load initializes the routing of the application.
func Load(cfg *config.Config) http.Handler {
	mux := chi.NewRouter()

	mux.Use(hlog.NewHandler(log.Logger))
	mux.Use(hlog.RemoteAddrHandler("ip"))
	mux.Use(hlog.URLHandler("path"))
	mux.Use(hlog.MethodHandler("method"))
	mux.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

	mux.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Debug().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))

	mux.Use(middleware.Timeout(MiddlewareTimeout))
	mux.Use(middleware.RealIP)
	mux.Use(metrics.DurationMetrics(&cfg.Metrics.Metrics))
	mux.Use(header.Version)
	mux.Use(header.Cache)
	mux.Use(header.Secure)
	mux.Use(header.Options)

	mux.Route(cfg.Server.Root, func(root chi.Router) {
		root.Get("/{code}.html", errorpagesHandler.NewHandler(cfg))
		root.Get("/healthz", healthHandler.NewHandler())

		if cfg.Server.Pprof {
			root.Mount("/debug", middleware.Profiler())
		}
	})

	mux.NotFound(notfound.NewHandler(cfg))

	return mux
}

// Metrics initializes the routing of metrics and health.
func Metrics(cfg *config.Config) http.Handler {
	mux := chi.NewRouter()

	mux.Use(hlog.NewHandler(log.Logger))
	mux.Use(hlog.RemoteAddrHandler("ip"))
	mux.Use(hlog.URLHandler("path"))
	mux.Use(hlog.MethodHandler("method"))
	mux.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

	mux.Use(middleware.Timeout(MiddlewareTimeout))
	mux.Use(middleware.RealIP)
	mux.Use(header.Version)
	mux.Use(header.Cache)
	mux.Use(header.Secure)
	mux.Use(header.Options)

	mux.Route("/", func(root chi.Router) {
		root.Get("/metrics", metricsHandler.NewHandler(cfg))
		root.Get("/healthz", healthHandler.NewHandler())
	})

	mux.NotFound(notfound.NewHandler(cfg))

	return mux
}

// Curves provides optionally a list of secure curves.
func Curves(cfg *config.Config) []tls.CurveID {
	if cfg.Server.StrictCurves {
		return []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		}
	}

	return nil
}

// Ciphers provides optionally a list of secure ciphers.
func Ciphers(cfg *config.Config) []uint16 {
	if cfg.Server.StrictCiphers {
		return []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		}
	}

	return nil
}
