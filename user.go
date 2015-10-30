package isumm

import (
	"html/template"
	"net/http"
	"os"

	"appengine"
	"appengine/user"
)

func IsUserAllowed(u *user.User) bool {
	if os.Getenv("RUN_WITH_DEVAPPSERVER") == "1" {
		return u.Email == AllowedTestUser
	}
	return u.Email == AllowedUser
}

var notAllowedTemplate = template.Must(template.ParseFiles("static/not_allowed.template.html"))

type notAllowedParams struct {
	User      string
	LogoutUrl string
}

func InvalidUserPage(c appengine.Context, w http.ResponseWriter, r *http.Request, u *user.User) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusForbidden)

	logoutUrl, err := LogoutURL(c, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := notAllowedTemplate.Execute(w, notAllowedParams{u.String(), logoutUrl}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
