package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.owncloud.com/owncloud-ops/errors/pkg/config"
)

//go:embed dist/*
var embeddedTemplates embed.FS

// Load initializes the template files.
func Load(cfg *config.Config) *template.Template {
	tpls := template.New("")

	err := fs.WalkDir(embeddedTemplates, ".", func(name string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() {
			return nil
		}

		if forbiddenExtension(filepath.Ext(dir.Name())) {
			return nil
		}

		content, err := fs.ReadFile(
			embeddedTemplates,
			name,
		)
		if err != nil {
			return fmt.Errorf("failed to read embedded template file: %w", err)
		}

		_, _ = tpls.New(
			strings.TrimPrefix(
				dir.Name(),
				"dist/",
			),
		).Parse(
			string(content),
		)

		return nil
	})
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to parse builtin templates")
	}

	//nolint:nestif
	if cfg.Server.Templates != "" {
		if stat, err := os.Stat(cfg.Server.Templates); os.IsNotExist(err) || !stat.IsDir() {
			log.Warn().
				Err(err).
				Msg("Custom templates directory does not exit")

			return tpls
		}

		err := filepath.Walk(cfg.Server.Templates, func(name string, dir fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if dir.IsDir() {
				return nil
			}

			if forbiddenExtension(filepath.Ext(dir.Name())) {
				return nil
			}

			content, err := os.ReadFile(
				name,
			)
			if err != nil {
				return fmt.Errorf("failed to read custom template file: %w", err)
			}

			_, _ = tpls.New(
				strings.TrimPrefix(
					strings.TrimPrefix(
						dir.Name(),
						cfg.Server.Templates,
					),
					"/",
				),
			).Parse(
				string(content),
			)

			return nil
		})
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Failed to parse custom templates")
		}
	}

	return tpls
}

func forbiddenExtension(ext string) bool {
	allowedExtensions := []string{
		".tmpl",
		".html",
		".json",
	}

	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return false
		}
	}

	return true
}
