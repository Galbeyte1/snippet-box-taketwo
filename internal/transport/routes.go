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

	// Dynamic (CSRF only, since Gorilla doesn’t have LoadAndSave)
	dynamic := alice.New(
		noSurf,
	)

	// Public dynamic
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected dynamic (auth on top)
	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// Standard across everything
	standard := alice.New(
		app.recoverPanic,
		app.logRequest,
		commonHeaders,
	)

	return standard.Then(mux)
}
