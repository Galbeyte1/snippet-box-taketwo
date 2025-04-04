package transport

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/config"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/helpers"
)

func Home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Go")

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

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			helpers.ServerError(app, err)
		}
	}
}

func SnippetView(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
		}

		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	}
}

func SnippetCreate(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Display a form for creating a new snippet..."))
	}
}

func SnippetCreatePost(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Save a new snippet..."))
	}
}
