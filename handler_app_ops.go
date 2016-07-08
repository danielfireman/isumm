package isumm

import (
	"net/http"
	"text/template"

	"appengine"
	"appengine/user"
)

var appOpsTemplate = template.Must(template.ParseFiles("static/ops.template.html"))

type appOpsParams struct {
	User        string
	LogoutURL   string
	Investments []*Investment
}

func AppOps(w http.ResponseWriter, r *http.Request) {
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

	investments, err := GetInvestments(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params := appParams{
		User:        u.String(),
		LogoutURL:   logoutUrl,
		Investments: investments}
	if err := appOpsTemplate.Execute(w, params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
