package isumm

import (
	"fmt"
	"net/http"

	"appengine"
)

const (
	InvParamKey  = "key"
	InvParamName = "name"
)

func Inv(w http.ResponseWriter, r *http.Request) {
	if err := handleInv(appengine.NewContext(r), r); err != nil {
		http.Error(w, err.Msg, err.Code)
		return
	}
	// Redirect is here because it leads to a panic (invalid memory address) when testing. It looks like a
	// appengine bug.
	http.Redirect(w, r, "/app", http.StatusFound)
}

func handleInv(c appengine.Context, r *http.Request) *handlingError {
	action := r.FormValue(ActionParam)
	switch action {
	case DeleteAction:
		if err := deleteInvestment(c, r.FormValue(InvParamKey)); err != nil {
			return &handlingError{err.Error(), http.StatusPreconditionFailed}
		}
	case PostAction:
		if err := postInvestment(c, r.FormValue(InvParamName), r.FormValue(InvParamKey)); err != nil {
			return &handlingError{err.Error(), http.StatusPreconditionFailed}
		}
	default:
		return &handlingError{fmt.Sprintf("Invalid action:\"%s\"", action), http.StatusBadRequest}
	}
	return nil
}

func deleteInvestment(c appengine.Context, k string) error {
	if k == "" {
		return fmt.Errorf("Investment key can not be empty.")
	}
	return DeleteInvestment(c, k)
}

func postInvestment(c appengine.Context, name, key string) error {
	if name == "" {
		return fmt.Errorf("Investment name can not be empty.")
	}
	if key == "" {
		return PutInvestment(c, &Investment{Name: name})
	}
	inv, err := GetInvestment(c, key)
	if err != nil {
		return fmt.Errorf("Failure fetching investment to be updated. Name:%s Key:%s", name, key)
	}
	inv.Name = name
	return PutInvestment(c, inv)
}
