package transport

import (
	"net/http"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/config"
)

func (app *Application) Routes(cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.StaticDir))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.protect(app.snippetCreate))
	mux.HandleFunc("POST /snippet/create", app.protect(app.snippetCreatePost))

	mux.HandleFunc("GET /user/signup", app.userSignup)
	mux.HandleFunc("POST /user/signup", app.userSignupPost)
	mux.HandleFunc("GET /user/login", app.userLogin)
	mux.HandleFunc("POST /user/login", app.userLoginPost)
	mux.HandleFunc("POST /user/logout", app.protect(app.userLogoutPost))

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
