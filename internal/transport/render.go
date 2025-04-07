package transport

import (
	"fmt"
	"net/http"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/templates"
)

func (app *Application) Render(w http.ResponseWriter, r *http.Request, status int, page string, data templates.TemplateData) {
	ts, ok := app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServerError(err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.ServerError(err)
	}
}
