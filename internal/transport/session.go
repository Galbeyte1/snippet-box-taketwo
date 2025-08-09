package transport

import "net/http"

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
