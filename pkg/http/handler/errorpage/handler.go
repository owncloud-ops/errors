package errorpages

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.owncloud.com/owncloud-ops/errors/pkg/config"
	"github.owncloud.com/owncloud-ops/errors/pkg/http/core"
)

// NewHandler creates handler for error pages serving.
func NewHandler(cfg *config.Config) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		// defer handleMetrics(time.Now(), req.ProtoMajor, req.ProtoMinor)
		core.SetClientFormat(writer, core.PlainTextContentType)

		if code, err := strconv.Atoi(chi.URLParam(req, "code")); err == nil {
			core.RespondWithErrorPage(req, writer, cfg, code)
		} else {
			code = http.StatusInternalServerError
			log.Error().
				Int("code", code).
				Msg("Invalid request code extracted from request")

			writer.WriteHeader(code)
			_, _ = io.WriteString(writer, "cannot extract requested code from the request")
		}
	}
}
