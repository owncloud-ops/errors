package core

import (
	"io"
	"net/http"

	"github.com/owncloud-ops/errors/pkg/config"
	"github.com/owncloud-ops/errors/pkg/errors"
	"github.com/owncloud-ops/errors/pkg/templates"
	"github.com/rs/zerolog/log"
)

// Payload represents the payload for template rendering.
type Payload struct {
	Status int
	Error  string
	Title  string
}

func RespondWithErrorPage(
	req *http.Request,
	writer http.ResponseWriter,
	cfg *config.Config,
	pageCode int,
) {
	errorTemplate := "html.tmpl"
	availableErrors := errors.Load(cfg)
	msg, ok := availableErrors[pageCode]

	if !ok {
		msg = http.StatusText(pageCode)
	}

	clientWant := ClientWantFormat(req)

	writer.Header().Set("X-Robots-Tag", "noindex") // block Search indexing
	SetClientFormat(writer, PlainTextContentType)  // set default content type

	switch {
	case clientWant == JSONContentType: // JSON
		{
			errorTemplate = "json.tmpl"
			SetClientFormat(writer, JSONContentType)
		}

	default: // HTML
		{
			SetClientFormat(writer, HTMLContentType)
		}
	}

	writer.WriteHeader(pageCode)

	if err := templates.Load(cfg).ExecuteTemplate(
		writer,
		errorTemplate,
		Payload{
			Status: pageCode,
			Error:  msg,
			Title:  cfg.Server.ErrorsTitle,
		},
	); err != nil {
		log.Error().
			Err(err).
			Str("template", errorTemplate).
			Msg("Failed to execute template")

		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(writer, "template "+errorTemplate+" not exists")
	}
}
