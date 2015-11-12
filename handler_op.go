package isumm

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"appengine"
)

const (
	OpsParamType  = "type"
	OpsParamValue = "value"
	OpsParamDate  = "date"
	OpsParamInv   = "inv"
)

type handlingError struct {
	Msg  string
	Code int
}

func Op(w http.ResponseWriter, r *http.Request) {
	if err := handleOp(appengine.NewContext(r), r); err != nil {
		http.Error(w, err.Msg, err.Code)
		return
	}
	// Redirect is here because it leads to a panic (invalid memory address) when testing. It looks like a
	// appengine bug.
	http.Redirect(w, r, "/app", http.StatusFound)
}

func handleOp(c appengine.Context, r *http.Request) *handlingError {
	// 1st phase: parameters extraction and validation.
	action := r.FormValue("action")
	invStr := r.FormValue(OpsParamInv)
	if invStr == "" {
		return &handlingError{"Investment key can not be empty.", http.StatusPreconditionFailed}
	}
	var (
		err   error
		op    Operation
		index int
	)
	switch action {
	case "d":
		index, err = strconv.Atoi(r.FormValue("index"))
		if err != nil {
			return &handlingError{fmt.Sprintf("Invalid operation index: %s", r.FormValue("index")), http.StatusPreconditionFailed}
		}
	default:
		op, err = NewOperationFromString(r.FormValue(OpsParamType), r.FormValue(OpsParamValue), r.FormValue(OpsParamDate))
		if err != nil {
			return &handlingError{err.Error(), http.StatusPreconditionFailed}
		}
	}
	// 2nd phase: updating investment.
	inv, err := GetInvestment(c, invStr)
	if err != nil {
		return &handlingError{err.Error(), http.StatusInternalServerError}
	}
	switch action {
	case "d":
		inv.Ops = append(inv.Ops[:index], inv.Ops[index+1:]...)

	default:
		inv.Ops = append(inv.Ops, op)
		sort.Sort(inv.Ops)
	}
	if err := PutInvestment(c, inv); err != nil {
		return &handlingError{err.Error(), http.StatusInternalServerError}
	}
	return nil
}
