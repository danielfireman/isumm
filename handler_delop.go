package isumm

import (
	"html/template"
	"net/http"

	"appengine"
)

var delOpTemplate = template.Must(template.ParseFiles("static/delop.template.html"))

type delOpParams struct {
	InvKey  string
	InvName string
	Ops     Operations
}

func DelOp(w http.ResponseWriter, r *http.Request) {
	invStr := r.FormValue("inv")
	if invStr == "" {
		http.Error(w, "Investment key can not be empty.", http.StatusPreconditionFailed)
		return
	}
	c := appengine.NewContext(r)
	inv, err := GetInvestment(c, invStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := delOpTemplate.Execute(w, delOpParams{inv.Key, inv.Name, inv.Ops}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
