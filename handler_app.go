package isumm

import (
	"net/http"
	"text/template"

	"appengine"
	"appengine/user"
)

var appTemplate = template.Must(template.ParseFiles("static/app.template"))

type appParams struct {
	User      string
	LogoutUrl string
}

func App(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c) // Login is mandatory on this page. No need to check nil value here.
	if !IsUserAllowed(u) {
		InvalidUserPage(c, w, r, u)
		return
	}
	logoutUrl, err := LogoutURL(c, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := appTemplate.Execute(w, appParams{u.String(), logoutUrl}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
