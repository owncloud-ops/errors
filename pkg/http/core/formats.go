package core

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type ContentType = byte

const (
	UnknownContentType ContentType = iota // should be first
	JSONContentType
	HTMLContentType
	PlainTextContentType
)

const (
	// FormatHeader name of the header used to extract the format.
	FormatHeader = "X-Format"
)

func ClientWantFormat(req *http.Request) ContentType {
	// parse "Content-Type" header (e.g.: `application/json;charset=UTF-8`)
	if ct := strings.ToLower(req.Header.Get("Content-type")); len(ct) > 4 { //nolint:gomnd
		return mimeTypeToContentType(ct)
	}

	// parse `X-Format` header (aka `Accept`) for the Ingress support
	// e.g.: `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8`
	if h := strings.ToLower(strings.TrimSpace(req.Header.Get(FormatHeader))); len(h) > 2 { //nolint:gomnd,nestif
		type format struct {
			mimeType string
			weight   float32
		}

		formats := make([]format, 0, 8) //nolint:gomnd

		for _, b := range strings.FieldsFunc(h, func(r rune) bool { return r == ',' }) {
			if idx := strings.Index(b, ";q="); idx > 0 && idx < len(b) {
				f := format{b[0:idx], 0}

				if len(b) > idx+3 {
					if weight, err := strconv.ParseFloat(b[idx+3:], 32); err == nil {
						f.weight = float32(weight)
					}
				}

				formats = append(formats, f)
			} else {
				formats = append(formats, format{b, 1})
			}
		}

		switch l := len(formats); {
		case l == 0:
			return UnknownContentType

		case l == 1:
			return mimeTypeToContentType(formats[0].mimeType)

		default:
			sort.SliceStable(formats, func(i, j int) bool { return formats[i].weight > formats[j].weight })

			return mimeTypeToContentType(formats[0].mimeType)
		}
	}

	return UnknownContentType
}

func SetClientFormat(writer http.ResponseWriter, t ContentType) {
	switch t {
	case JSONContentType:
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	case HTMLContentType:
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	case PlainTextContentType:
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

func mimeTypeToContentType(mimeType string) ContentType {
	switch {
	case strings.Contains(mimeType, "application/json"), strings.Contains(mimeType, "text/json"):
		return JSONContentType

	case strings.Contains(mimeType, "text/html"):
		return HTMLContentType

	case strings.Contains(mimeType, "text/plain"):
		return PlainTextContentType
	}

	return UnknownContentType
}
