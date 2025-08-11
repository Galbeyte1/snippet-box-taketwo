package transport

import "net/http"

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

func (app *Application) RenewSession(w http.ResponseWriter, r *http.Request, userID int) error {
	if old, err := app.SessionStore.Get(r, "session"); err == nil {
		old.Options.MaxAge = -1
		if err := old.Save(r, w); err != nil {
			return err
		}
	}

	sess, err := app.SessionStore.New(r, "session")
	if err != nil {
		return err
	}
	sess.Values["authenticatedUserID"] = userID
	return sess.Save(r, w)
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
