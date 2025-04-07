package transport

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.Snippets.Latest()
	if err != nil {
		app.ServerError(err)
		return
	}
	app.Render(w, r, http.StatusOK, "home.tmpl", templates.TemplateData{
		Snippets: snippets,
	})
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
			app.ServerError(err)
		}
		return
	}

	app.Render(w, r, http.StatusOK, "view.tmpl", templates.TemplateData{
		Snippet: snippet,
	})
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := `الإمام الشافعي`
	content := `دع الأيام تفعل ما تشاءُ وطب نفساً إذا حكم القضاءُ ولا تجزع لحادثة الليالي فما لحوادث الدنيا بقاءُ`
	expires := 20

	id, err := app.Snippets.Insert(title, content, expires)
	if err != nil {
		app.ServerError(err)
		return
	}
	log.Println("Successfully inserted snippet with ID", id) // <-- And this
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
