package transport

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/config"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/helpers"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/models"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/templates"
)

/*
File Server and Template Parsing example

Step	Browser action												Go server reaction
1 		Browser requests http://localhost:8080/static/style.css		FileServer looks inside ./static/style.css and streams it
2		Browser requests http://localhost:8080/						HomeHandler runs, tmpl.Execute renders index.html with Name = "Alice"
3		index.html links to static/style.css						Browser automatically requests style.css, served again by FileServer

Building a server cache for templates WHILE file server caching handled by os/browser
*/

func home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Go")

		snippets, err := app.Snippets.Latest()
		if err != nil {
			helpers.ServerError(app, err)
			return
		}

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/home.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			helpers.ServerError(app, err)
			return
		}

		data := templates.TemplateData{
			Snippets: snippets,
		}

		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			helpers.ServerError(app, err)
		}
	}
}

func snippetView(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
		}
		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNotRecord) {
				http.NotFound(w, r)
			} else {
				helpers.ServerError(app, err)
			}
			return
		}

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/view.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			helpers.ServerError(app, err)
			return
		}

		data := templates.TemplateData{
			Snippet: snippet,
		}

		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			helpers.ServerError(app, err)
		}
	}

}

func snippetCreate(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Display a form for creating a new snippet..."))
	}
}

func snippetCreatePost(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := `الإمام الشافعي`
		content := `دع الأيام تفعل ما تشاءُ وطب نفساً إذا حكم القضاءُ ولا تجزع لحادثة الليالي فما لحوادث الدنيا بقاءُ`
		expires := 20

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			helpers.ServerError(app, err)
			return
		}
		log.Println("Successfully inserted snippet with ID", id) // <-- And this
		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}
