package isumm

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"appengine/aetest"
)

const (
	typeStr        = "1"
	valueStr       = "2.23"
	dateStr        = "2015-12-30"
	firstOpDateStr = "2015-12-31"
)

func TestOpHandler_AddSucess(t *testing.T) {
	c, err := aetest.NewContext(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	// The test has two operations. One is added at investment creation and the other is added via handler.
	firstOp, _ := NewOperationFromString(typeStr, valueStr, firstOpDateStr)
	secondOp, _ := NewOperationFromString(typeStr, valueStr, dateStr)
	// Adding investment along with first operation.
	inv := &Investment{Name: "testInv", Ops: Operations{firstOp}}
	if err := PutInvestment(c, inv); err != nil {
		t.Fatal(err)
	}
	// This should add an operation that is exactly like second operation.
	form := url.Values{OpsParamInv: {inv.Key}, OpsParamType: {typeStr}, OpsParamValue: {valueStr}, OpsParamDate: {dateStr}}
	if err := handleOp(c, &http.Request{Form: form}); err != nil {
		t.Fatalf("handleOp want nil got:%q", err)
	}
	gotInv, err := GetInvestment(c, inv.Key)
	if err != nil {
		t.Errorf("getInvestment returned an error: %q", err)
	}
	// Even though firstOp happened latter, operations must be chronologically sorted.
	if !reflect.DeepEqual(firstOp, gotInv.Ops[0]) { // There is already one operation inserted.
		t.Errorf("want:%v got:%v", firstOp, gotInv.Ops[0])
	}
	if !reflect.DeepEqual(secondOp, gotInv.Ops[1]) { // There is already one operation inserted.
		t.Errorf("want:%v got:%v", secondOp, gotInv.Ops[1])
	}
}

func TestOpHandler_AddFailure(t *testing.T) {
	c, err := aetest.NewContext(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	op, _ := NewOperationFromString(typeStr, valueStr, dateStr)
	inv := &Investment{Name: "testInv", Ops: Operations{op}}
	if err := PutInvestment(c, inv); err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		desc, k, t, v, d string
		action           string
		statusCode       int
	}{
		{
			desc: "invalid key",
			k:    "", t: typeStr, v: valueStr, d: dateStr,
			statusCode: http.StatusPreconditionFailed,
		},
		{
			desc: "invalid date",
			k:    inv.Key, t: typeStr, v: valueStr, d: "",
			statusCode: http.StatusPreconditionFailed,
		},
		{
			desc: "invalid investment",
			k:    "1", t: typeStr, v: valueStr, d: dateStr,
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, testCase := range testCases {
		form := url.Values{OpsParamInv: {testCase.k}, OpsParamType: {testCase.t}, OpsParamValue: {testCase.v}, OpsParamDate: {testCase.d}}
		gotErr := handleOp(c, &http.Request{Form: form})
		if gotErr == nil {
			t.Errorf("want error got nil. Test case: %+v", testCase.desc)
		}
		if testCase.statusCode != gotErr.Code {
			t.Errorf("(%s) status code want:%d got:%d", testCase.desc, testCase.statusCode, gotErr.Code)
		}
	}
}
