package command

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.owncloud.com/owncloud-ops/errors/pkg/http/router"
	"github.owncloud.com/owncloud-ops/errors/pkg/metrics"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start integrated server",
	Run:   serverAction,
}

const (
	defaultMetricsAddr         = "0.0.0.0:8081"
	defaultServerAddr          = "0.0.0.0:8080"
	defaultServerPprof         = false
	defaultServerRoot          = "/"
	defaultServerHost          = "http://localhost:8080"
	defaultServerCert          = ""
	defaultServerKey           = ""
	defaultServerStrictCurves  = false
	defaultServerStrictCiphers = false
	defaultServerTemplates     = ""
	defaultServerErrors        = ""
	defaultServerErrorsTitle   = ""
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().String("metrics-addr", defaultMetricsAddr, "Address to bind the metrics")
	viper.SetDefault("metrics.addr", defaultMetricsAddr)
	_ = viper.BindPFlag("metrics.addr", serverCmd.PersistentFlags().Lookup("metrics-addr"))

	serverCmd.PersistentFlags().String("metrics-token", "", "Token to make metrics secure")
	viper.SetDefault("metrics.token", "")
	_ = viper.BindPFlag("metrics.token", serverCmd.PersistentFlags().Lookup("metrics-token"))

	serverCmd.PersistentFlags().String("server-addr", defaultServerAddr, "Address to bind the server")
	viper.SetDefault("server.addr", defaultServerAddr)
	_ = viper.BindPFlag("server.addr", serverCmd.PersistentFlags().Lookup("server-addr"))

	serverCmd.PersistentFlags().Bool("server-pprof", defaultServerPprof, "Enable pprof debugging")
	viper.SetDefault("server.pprof", defaultServerPprof)
	_ = viper.BindPFlag("server.pprof", serverCmd.PersistentFlags().Lookup("server-pprof"))

	serverCmd.PersistentFlags().String("server-root", defaultServerRoot, "Root path of the server")
	viper.SetDefault("server.root", defaultServerRoot)
	_ = viper.BindPFlag("server.root", serverCmd.PersistentFlags().Lookup("server-root"))

	serverCmd.PersistentFlags().String("server-host", defaultServerHost, "External access to server")
	viper.SetDefault("server.host", defaultServerHost)
	_ = viper.BindPFlag("server.host", serverCmd.PersistentFlags().Lookup("server-host"))

	serverCmd.PersistentFlags().String("server-cert", defaultServerCert, "Path to cert for SSL encryption")
	viper.SetDefault("server.cert", defaultServerCert)
	_ = viper.BindPFlag("server.cert", serverCmd.PersistentFlags().Lookup("server-cert"))

	serverCmd.PersistentFlags().String("server-key", defaultServerKey, "Path to key for SSL encryption")
	viper.SetDefault("server.key", defaultServerKey)
	_ = viper.BindPFlag("server.key", serverCmd.PersistentFlags().Lookup("server-key"))

	serverCmd.PersistentFlags().Bool("strict-curves", defaultServerStrictCurves, "Use strict SSL curves")
	viper.SetDefault("server.strict_curves", defaultServerStrictCurves)
	_ = viper.BindPFlag("server.strict_curves", serverCmd.PersistentFlags().Lookup("strict-curves"))

	serverCmd.PersistentFlags().Bool("strict-ciphers", defaultServerStrictCiphers, "Use strict SSL ciphers")
	viper.SetDefault("server.strict_ciphers", defaultServerStrictCiphers)
	_ = viper.BindPFlag("server.strict_ciphers", serverCmd.PersistentFlags().Lookup("strict-ciphers"))

	serverCmd.PersistentFlags().String("templates-path", defaultServerTemplates, "Path for overriding templates")
	viper.SetDefault("server.templates", defaultServerTemplates)
	_ = viper.BindPFlag("server.templates", serverCmd.PersistentFlags().Lookup("templates-path"))

	serverCmd.PersistentFlags().String("errors-path", defaultServerErrors, "Path for overriding errors")
	viper.SetDefault("server.errors", defaultServerErrors)
	_ = viper.BindPFlag("server.errors", serverCmd.PersistentFlags().Lookup("errors-path"))

	serverCmd.PersistentFlags().String("errors-title", defaultServerErrorsTitle, "String for overriding errors title")
	viper.SetDefault("server.errors_title", defaultServerErrorsTitle)
	_ = viper.BindPFlag("server.errors_title", serverCmd.PersistentFlags().Lookup("errors-title"))
}

//nolint:revive
func serverAction(ccmd *cobra.Command, args []string) {
	const (
		RunTimeout       = 3 * time.Second
		HTTPReadTimeout  = 5 * time.Second
		HTTPWriteTimeout = 10 * time.Second
	)

	var group run.Group

	//nolint:nestif
	if cfg.Server.Cert != "" && cfg.Server.Key != "" {
		cert, err := tls.LoadX509KeyPair(
			cfg.Server.Cert,
			cfg.Server.Key,
		)
		if err != nil {
			log.Info().
				Err(err).
				Msg("Failed to load certificates")

			os.Exit(1)
		}

		server := &http.Server{
			Addr:         cfg.Server.Addr,
			Handler:      router.Load(cfg),
			ReadTimeout:  HTTPReadTimeout,
			WriteTimeout: HTTPWriteTimeout,
			TLSConfig: &tls.Config{
				PreferServerCipherSuites: true,
				MinVersion:               tls.VersionTLS12,
				CurvePreferences:         router.Curves(cfg),
				CipherSuites:             router.Ciphers(cfg),
				Certificates:             []tls.Certificate{cert},
			},
		}

		group.Add(func() error {
			log.Info().
				Str("addr", cfg.Server.Addr).
				Msg("Starting HTTPS server")

			if err := server.ListenAndServeTLS("", ""); err != nil {
				return fmt.Errorf("failed to start https server: %w", err)
			}

			return nil
		}, func(reason error) {
			ctx, cancel := context.WithTimeout(context.Background(), RunTimeout)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				log.Error().
					Err(err).
					Msg("Failed to shutdown HTTPS gracefully")

				return
			}

			log.Info().
				Err(reason).
				Msg("Shutdown HTTPS gracefully")
		})
	} else {
		server := &http.Server{
			Addr:         cfg.Server.Addr,
			Handler:      router.Load(cfg),
			ReadTimeout:  HTTPReadTimeout,
			WriteTimeout: HTTPWriteTimeout,
		}

		group.Add(func() error {
			log.Info().
				Str("addr", cfg.Server.Addr).
				Msg("Starting HTTP server")

			if err := server.ListenAndServe(); err != nil {
				return fmt.Errorf("failed to start http server: %w", err)
			}

			return nil
		}, func(reason error) {
			ctx, cancel := context.WithTimeout(context.Background(), RunTimeout)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				log.Error().
					Err(err).
					Msg("Failed to shutdown HTTP gracefully")

				return
			}

			log.Info().
				Err(reason).
				Msg("Shutdown HTTP gracefully")
		})
	}

	{
		cfg.Metrics.Reg, cfg.Metrics.Metrics = metrics.NewRegistry(), metrics.NewMetrics()
		server := &http.Server{
			Addr:         cfg.Metrics.Addr,
			Handler:      router.Metrics(cfg),
			ReadTimeout:  HTTPReadTimeout,
			WriteTimeout: HTTPWriteTimeout,
		}

		group.Add(func() error {
			log.Info().
				Str("addr", cfg.Metrics.Addr).
				Msg("Starting metrics server")

			if err := cfg.Metrics.Metrics.Register(cfg.Metrics.Reg); err != nil {
				return fmt.Errorf("failed register metrics: %w", err)
			}

			if err := server.ListenAndServe(); err != nil {
				return fmt.Errorf("failed to start metrics server: %w", err)
			}

			return nil
		}, func(reason error) {
			ctx, cancel := context.WithTimeout(context.Background(), RunTimeout)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				log.Error().
					Err(err).
					Msg("Failed to shutdown metrics gracefully")

				return
			}

			log.Info().
				Err(reason).
				Msg("Shutdown metrics gracefully")
		})
	}

	{
		stop := make(chan os.Signal, 1)

		group.Add(func() error {
			signal.Notify(stop, os.Interrupt)

			<-stop

			return nil
		}, func(err error) {
			close(stop)
		})
	}

	if err := group.Run(); err != nil {
		os.Exit(1)
	}
}
