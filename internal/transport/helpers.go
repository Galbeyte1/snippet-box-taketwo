package transport

import (
	"errors"
	"net/http"

	"github.com/go-playground/form/v4"
)

const (
	sessionName      = "session"
	sessionUserIDKey = "authenticatedUserID"
)

func (app *Application) PopFlash(w http.ResponseWriter, r *http.Request) string {
	sess, err := app.SessionStore.Get(r, "flash")
	if err != nil {
		return ""
	}

	flashes := sess.Flashes()
	if len(flashes) == 0 {
		return ""
	}

	_ = sess.Save(r, w)
	if msg, ok := flashes[0].(string); ok {
		return msg
	}
	return ""
}

func (app *Application) DecodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return nil
	}

	err = app.FormDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *Application) isAuthenticated(r *http.Request) bool {
	sess, err := app.SessionStore.Get(r, sessionName)
	if err != nil {
		return false
	}

	v, ok := sess.Values[sessionUserIDKey]
	if !ok {
		return false
	}
	switch x := v.(type) {
	case int:
		return x > 0
	case int64:
		return x > 0
	case string:
		return x != "" && x != "0"
	default:
		return false
	}
}

func (app *Application) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		// trace = string(debug.Stack())
	)

	app.Logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
