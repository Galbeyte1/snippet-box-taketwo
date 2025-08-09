package transport

import "net/http"

func (app *Application) RenewSession(w http.ResponseWriter, r *http.Request, userID int) error {
	old, _ := app.SessionStore.Get(r, "session")
	old.Options.MaxAge = -1
	if err := old.Save(r, w); err != nil {
		return err
	}

	sess, _ := app.SessionStore.New(r, "session")
	sess.Values["authenticatedUserID"] = userID
	return sess.Save(r, w)
}
