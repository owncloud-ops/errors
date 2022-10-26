package handler

import (
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.owncloud.com/owncloud-ops/errors/pkg/config"
	"github.owncloud.com/owncloud-ops/errors/pkg/errors"
	"github.owncloud.com/owncloud-ops/errors/pkg/model"
	"github.owncloud.com/owncloud-ops/errors/pkg/templates"
)

// General is used to handle all error pages.
func General(cfg *config.Config) http.HandlerFunc {
	availableErrors := errors.Load(cfg)

	return func(w http.ResponseWriter, req *http.Request) {
		defer handleMetrics(time.Now(), req.ProtoMajor, req.ProtoMinor)

		code := detectCode(req)
		format := detectFormat(req)
		file := parseFormat(format)

		if http.StatusText(code) == "" {
			log.Info().
				Int("code", code).
				Msg("Invalid request code extracted from request")
			code = http.StatusInternalServerError
		}

		w.Header().Set("Content-Type", format)
		w.WriteHeader(code)

		msg, ok := availableErrors[code]

		if !ok {
			msg = http.StatusText(code)
		}

		if err := templates.Load(cfg).ExecuteTemplate(
			w,
			file,
			model.Payload{
				Status: code,
				Error:  msg,
				Title:  cfg.Server.ErrorsTitle,
			},
		); err != nil {
			log.Error().
				Err(err).
				Str("template", file).
				Msg("Failed to execute template")

			_, _ = io.WriteString(w, http.StatusText(code))
		}
	}
}

func detectCode(req *http.Request) int {
	if val := req.Header.Get("X-Code"); val != "" {
		code, err := strconv.Atoi(val)
		if err != nil {
			code = 404

			log.Info().
				Err(err).
				Int("code", code).
				Msg("Failed to parse code")
		}

		return code
	}

	name := path.Base(
		req.URL.String(),
	)

	if name != "/" {
		val := strings.TrimSuffix(
			name,
			filepath.Ext(name),
		)

		code, err := strconv.Atoi(val)
		if err != nil {
			code = 404

			log.Info().
				Err(err).
				Str("code", val).
				Msg("Failed to parse path")
		}

		return code
	}

	return 404
}

func detectFormat(req *http.Request) string {
	if val := req.Header.Get("X-Format"); val != "" {
		return val
	}

	name := path.Base(
		req.URL.String(),
	)

	if name != "/" {
		switch filepath.Ext(name) {
		case ".html":
			return "text/html"
		case ".json":
			return "application/json"
		default:
			log.Info().
				Str("format", name).
				Msg("Failed to detect format")

			return "text/html"
		}
	}

	format := "text/html"

	log.Info().
		Str("format", format).
		Msg("Format not specified, using default")

	return format
}

func parseFormat(format string) string {
	mediaType, _, _ := mime.ParseMediaType(format)
	cext, err := mime.ExtensionsByType(mediaType)
	if err != nil {
		log.Error().
			Err(err).
			Str("format", format).
			Msg("Failed to parse media type extension")

		return "html.tmpl"
	}

	if len(cext) == 0 {
		log.Info().
			Str("format", format).
			Msg("Could not detect media type extension")

		return "html.tmpl"
	}

	if cext[0] == ".htm" {
		return "html.tmpl"
	}

	return strings.TrimPrefix(
		cext[0],
		".",
	) + ".tmpl"
}
