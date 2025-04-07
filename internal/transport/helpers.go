package transport

import (
	"net/http"
)

func (app *Application) ServerError(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			method = r.Method
			uri    = r.URL.RequestURI()
			// trace = string(debug.Stack())
		)

		app.Logger.Error(err.Error(), "method", method, "uri", uri)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
