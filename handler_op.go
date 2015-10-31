package isumm

import (
	"fmt"
	"net/http"
	"strconv"
)

func Op(w http.ResponseWriter, r *http.Request) {
	strValue := r.FormValue("value")
	if strValue == "" {
		http.Error(w, "Value can not be empty.", http.StatusPreconditionFailed)
		return
	}
	value, err := strconv.ParseFloat(strValue, 32)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid value: %q", err), http.StatusPreconditionFailed)
		return
	}

	inv := r.FormValue("inv")
	if inv == "" {
		http.Error(w, "Please select an investment.", http.StatusPreconditionFailed)
		return
	}
	fmt.Fprintf(w, "inv:%s\n", inv)
	fmt.Fprintf(w, "value:%f\n", value)
	fmt.Fprintf(w, "type:%s\n", r.FormValue("type"))
}
