package transport

import (
	"net/http"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/config"
	"github.com/justinas/alice"
)

func (app *Application) Routes(cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.StaticDir))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.SessionManager.LoadAndSave)

	mux.HandleFunc("GET /{$}", dynamic.ThenFunc(app.home))
	mux.HandleFunc("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.HandleFunc("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.HandleFunc("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
