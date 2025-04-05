package config

import (
	"log/slog"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/models"
)

type Application struct {
	Logger   *slog.Logger
	Snippets *models.SnippetModel
}

type Config struct {
	Addr      string
	StaticDir string
	DSN       string
}
