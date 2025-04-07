package transport

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/templates"
)

func (app *Application) Render(w http.ResponseWriter, r *http.Request, status int, page string, data templates.TemplateData) {
	ts, ok := app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServerError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}
