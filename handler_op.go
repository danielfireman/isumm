package isumm

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"appengine"
)

func Op(w http.ResponseWriter, r *http.Request) {
	// Data validation.
	// TODO(danielfireman): Consider moving to Operation constructor.
	strValue := r.FormValue("value")
	if strValue == "" {
		http.Error(w, "Value can not be empty.", http.StatusPreconditionFailed)
		return
	}
	value, err := strconv.ParseFloat(strValue, 32)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid value: %s", strValue), http.StatusPreconditionFailed)
		return
	}
	inv := r.FormValue("inv")
	if inv == "" {
		http.Error(w, "Please select an investment.", http.StatusPreconditionFailed)
		return
	}
	strOpType := r.FormValue("type")
	opType, err := strconv.Atoi(strOpType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid operation type: %s", strOpType), http.StatusPreconditionFailed)
		return
	}
	strDate := r.FormValue("date")
	date, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid operation date: %s", strDate), http.StatusPreconditionFailed)
		return
	}
	// Finally putting everything together.
	c := appengine.NewContext(r)
	i, err := GetInvestment(c, inv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	op := Operation{
		Date:  date,
		Value: float32(value),
		Type:  OpType(opType),
	}
	i.Ops = append(i.Ops, op)
	if err := PutInvestment(c, i); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/app", http.StatusFound)
}
