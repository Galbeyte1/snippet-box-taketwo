package transport

import (
	"html/template"
	"log/slog"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/models"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/sessions"
)

type Application struct {
	Logger        *slog.Logger
	Snippets      *models.SnippetModel
	Users         *models.UserModel
	TemplateCache map[string]*template.Template
	FormDecoder   *form.Decoder
	SessionStore  sessions.Store
}
