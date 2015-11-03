package isumm

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"appengine"
)

func Op(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	action := r.FormValue("action")
	invStr := r.FormValue("inv")
	if invStr == "" {
		http.Error(w, "Investment key can not be empty.", http.StatusPreconditionFailed)
		return
	}
	inv, err := GetInvestment(c, invStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch action {
	case "d":
		indexStr := r.FormValue("index")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid operation index: %s", indexStr), http.StatusPreconditionFailed)
			return
		}
		inv.Ops = append(inv.Ops[:index], inv.Ops[index+1:]...)

	default:
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
		inv.Ops = append(inv.Ops, Operation{Date: date, Value: float32(value), Type: OpType(opType)})
	}
	if err := PutInvestment(c, inv); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/app", http.StatusFound)
}
