package transport

import (
	"html/template"
	"log/slog"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/models"
)

type Application struct {
	Logger        *slog.Logger
	Snippets      *models.SnippetModel
	TemplateCache map[string]*template.Template
}
