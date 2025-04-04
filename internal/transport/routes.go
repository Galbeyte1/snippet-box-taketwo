package transport

import (
	"net/http"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/config"
)

func Routes(app *config.Application, cfg config.Config) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.StaticDir))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home(app))
	mux.HandleFunc("GET /snippet/view/{id}", snippetView(app))
	mux.HandleFunc("GET /snippet/create", snippetCreate(app))
	mux.HandleFunc("POST /snippet/create", snippetCreatePost(app))

	return mux
}
