package isumm

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"appengine/aetest"
)

func TestHandleInv_Sucess(t *testing.T) {
	c, err := aetest.NewContext(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// Adding investment.
	name := "foo"
	form := url.Values{ActionParam: {PostAction}, InvParamKey: {""}, InvParamName: {name}}
	if err := handleInv(c, &http.Request{Form: form}); err != nil {
		t.Fatalf("handleOp want nil got:%q", err)
	}
	invs, err := GetInvestments(c)
	if err != nil {
		t.Fatalf("getInvestment returned an error: %q", err)
	}
	if len(invs) != 1 {
		t.Fatalf("#investments want:1 got:%d", len(invs))
	}
	inv := invs[0]
	if name != inv.Name {
		t.Fatalf("want:%v got:%v", name, inv.Name)
	}

	// Adding some ops just to be sure we are not wiping something out.
	inv.Ops = Operations{NewOperation(Balance, 1.2, time.Now())}
	if err := PutInvestment(c, inv); err != nil {
		t.Fatalf("putInvestment returned an error: %q", err)
	}

	// Renaming the recently added investment.
	newName := "foo"
	form = url.Values{ActionParam: {PostAction}, InvParamKey: {inv.Key}, InvParamName: {newName}}
	if err := handleInv(c, &http.Request{Form: form}); err != nil {
		t.Fatalf("handleOp want nil got:%q", err)
	}
	updatedInv, err := GetInvestment(c, inv.Key)
	if err != nil {
		t.Fatalf("getInvestment returned an error: %q", err)
	}
	if newName != updatedInv.Name {
		t.Fatalf("want:%v got:%v", name, inv.Name)
	}
	if len(updatedInv.Ops) != 1 {
		t.Fatalf("len(updatedInv.Ops) want:1 got:%d", len(updatedInv.Ops))
	}
	if !reflect.DeepEqual(inv.Ops, updatedInv.Ops) {
		t.Fatalf("ops want:%v got:%v", inv.Ops, updatedInv.Ops)
	}
}
