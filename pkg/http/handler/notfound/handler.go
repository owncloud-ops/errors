package notfound

import (
	"net/http"

	"github.com/owncloud-ops/errors/pkg/config"
	"github.com/owncloud-ops/errors/pkg/http/core"
)

// NewHandler creates handler missing requests handling.
func NewHandler(cfg *config.Config) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		core.RespondWithErrorPage(req, writer, cfg, http.StatusNotFound)
	}
}
