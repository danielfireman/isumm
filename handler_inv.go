package isumm

import (
	"net/http"

	"appengine"
)

func Inv(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	action := r.FormValue("action")
	switch action {
	case "d":
		k := r.FormValue("key")
		if k == "" {
			http.Error(w, "Investment key can not be empty.", http.StatusPreconditionFailed)
			return
		}
		if err := DeleteInvestment(c, k); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "Investment name can not be empty.", http.StatusPreconditionFailed)
			return
		}
		i := &Investment{Name: name, Key: r.FormValue("key")}
		if err := PutInvestment(c, i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/app", http.StatusFound)
}
