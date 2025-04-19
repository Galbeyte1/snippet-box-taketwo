package transport

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/models"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/templates"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/validator"
)

/*
File Server and Template Parsing example

Step	Browser action												Go server reaction
1 		Browser requests http://localhost:8080/static/style.css		FileServer looks inside ./static/style.css and streams it
2		Browser requests http://localhost:8080/						HomeHandler runs, tmpl.Execute renders index.html with Name = "Alice"
3		index.html links to static/style.css						Browser automatically requests style.css, served again by FileServer

Building a server cache for templates WHILE file server caching handled by os/browser
*/

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.Snippets.Latest()
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
	data := templates.NewTemplateData(r)
	data.Snippets = snippets
	app.Render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
	}
	snippet, err := app.Snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNotRecord) {
			http.NotFound(w, r)
		} else {
			app.ServerError(w, r, err)
		}
		return
	}

	data := templates.NewTemplateData(r)
	data.Snippet = snippet

	app.Render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := templates.NewTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.Render(w, r, http.StatusOK, "create.tmpl", data)

}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	var form snippetCreateForm

	err = app.FormDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must be equal 1, 7, or 365")

	if !form.Valid() {
		data := templates.NewTemplateData(r)
		data.Form = form
		app.Render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
