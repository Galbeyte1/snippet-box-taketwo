package transport

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/models"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/templates"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/validator"
)

/*
File Server and Template Parsing example

Step	Browser action												Go server reaction
1 		Browser requests http://localhost:8080/static/style.css		FileServer looks inside ./static/style.css and streams it
2		Browser requests http://localhost:8080/						HomeHandler runs, tmpl.Execute renders index.html with Name = "Alice"
3		index.html links to static/style.css						Browser automatically requests style.css, served again by FileServer

Building a server cache for templates WHILE file server caching handled by os/browser
*/

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.Snippets.Latest()
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
	data := templates.NewTemplateData(r)
	data.IsAuthenticated = app.isAuthenticated(r)
	data.Flash = app.PopFlash(w, r)
	data.Snippets = snippets
	app.Render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	snippet, err := app.Snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNotRecord) {
			http.NotFound(w, r)
		} else {
			app.ServerError(w, r, err)
		}
		return
	}

	data := templates.NewTemplateData(r)
	data.IsAuthenticated = app.isAuthenticated(r)
	data.Snippet = snippet
	data.Flash = app.PopFlash(w, r)

	app.Render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := templates.NewTemplateData(r)
	data.IsAuthenticated = app.isAuthenticated(r)
	data.Flash = app.PopFlash(w, r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.Render(w, r, http.StatusOK, "create.tmpl", data)

}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	var form snippetCreateForm

	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must be equal 1, 7, or 365")

	if !form.Valid() {
		data := templates.NewTemplateData(r)
		data.IsAuthenticated = app.isAuthenticated(r)
		data.Form = form
		app.Render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	session, err := app.SessionStore.Get(r, "flash")
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	session.AddFlash("Snippet successfully created!")

	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := templates.NewTemplateData(r)
	data.IsAuthenticated = app.isAuthenticated(r)
	data.Flash = app.PopFlash(w, r)
	data.Form = userSignupForm{}
	app.Render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be atleast 8 characters long")

	if !form.Valid() {
		data := templates.NewTemplateData(r)
		data.IsAuthenticated = app.isAuthenticated(r)
		data.Form = form
		app.Render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.Users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := templates.NewTemplateData(r)
			data.IsAuthenticated = app.isAuthenticated(r)
			data.Form = form
			app.Render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.ServerError(w, r, err)
		}

		return
	}

	session, err := app.SessionStore.Get(r, "flash")
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	session.AddFlash("Your signup was successful. Please log in")

	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	// This will serve/display the Login page template login.tmpl
	data := templates.NewTemplateData(r)
	data.IsAuthenticated = app.isAuthenticated(r)
	data.Flash = app.PopFlash(w, r)
	data.Form = userLoginForm{}
	app.Render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		// Serve up the page again
		data := templates.NewTemplateData(r)
		data.IsAuthenticated = app.isAuthenticated(r)
		data.Form = form
		app.Render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.Users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			// Serve up the page again
			data := templates.NewTemplateData(r)
			data.IsAuthenticated = app.isAuthenticated(r)
			data.Form = form
			app.Render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.ServerError(w, r, err)
		}
		return
	}

	if err := app.RenewSession(w, r, id); err != nil {
		app.ServerError(w, r, err)
		return
	}

	// flash message upon login
	if flash, err := app.SessionStore.Get(r, "flash"); err == nil {
		flash.AddFlash("Welcome back!")
		_ = flash.Save(r, w)
	}

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)

}

func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	sess, _ := app.SessionStore.Get(r, "session")
	sess.Options.MaxAge = -1
	_ = sess.Save(r, w)

	flash, _ := app.SessionStore.Get(r, "flash")
	flash.AddFlash("You've been logged out successfully!")
	_ = flash.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
