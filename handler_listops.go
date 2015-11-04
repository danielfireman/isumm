package isumm

import (
	"html/template"
	"net/http"

	"appengine"
)

var listOpsTemplate = template.Must(template.ParseFiles("static/listops.template.html"))

type listOpsParams struct {
	InvKey  string
	InvName string
	Ops     Operations
}

func ListOps(w http.ResponseWriter, r *http.Request) {
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
	if err := listOpsTemplate.Execute(w, listOpsParams{inv.Key, inv.Name, inv.Ops}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
