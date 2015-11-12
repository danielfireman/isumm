package isumm

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"
)

var (
	now       = time.Now()
	nextMonth = now.Add(32 * 24 * time.Hour)
	incMin    = now.Add(10 * time.Minute)
	opType    = "1"
	value     = "12"
	date      = "2006-01-02"
)

func TestSort(t *testing.T) {
	o := Operations{NewOperation(Balance, 1, now), NewOperation(Balance, 1, incMin)}
	sort.Sort(o)
	if o[0].Date().Before(o[1].Date()) {
		t.Errorf("want %s before than %s", o[0], o[1])
	}
}

func TestNewOperationFromString(t *testing.T) {
	validOp, _ := NewOperationFromString(opType, value, date)
	testCases := []struct {
		desc, t, v, d string
		want          *Operation
	}{
		// Invalid cases.
		{desc: "invalid value - empty", t: opType, v: "", d: date},
		{desc: "invalid value - chars", t: opType, v: "acb", d: date},
		{desc: "invalid op - empty string", t: "", v: value, d: date},
		{desc: "invalid op - type does not exist", t: "94879138ddffg", v: value, d: date},
		{desc: "invalid date", t: opType, v: value, d: "31/31/31"},
		// Valid cases.
		{desc: "valid - contain spaces", t: opType, v: fmt.Sprintf(" %s ", value), d: date, want: &validOp},
		{desc: "valid - perfect", t: opType, v: value, d: date, want: &validOp},
	}
	for _, test := range testCases {
		got, err := NewOperationFromString(test.t, test.v, test.d)
		switch {
		case test.want == nil:
			if err == nil {
				t.Errorf("got:nil expected:err. Test:(%+v)", test)
			}
		default:
			if !reflect.DeepEqual(test.want, &got) {
				t.Errorf("got:(%+v) want:(%+v). Test:(%+v)", got, test.want, test)
			}
		}
	}
}

func TestSummarize(t *testing.T) {
	data := []struct {
		ops  Operations
		summ []Summary
	}{
		{ // Simple case: One entry containing only balance.
			ops:  Operations{NewOperation(Balance, 1.2, now)},
			summ: []Summary{{Date: monthYear(now), Balance: 1.2, Change: 0}},
		},
		{ // Empty case.
			ops:  Operations{},
			summ: []Summary{},
		},
		{
			ops: Operations{
				NewOperation(Balance, 1.2, incMin),
				NewOperation(Deposit, 1.0, now),
				NewOperation(Balance, 2.2, nextMonth),
			},
			summ: []Summary{
				{Date: monthYear(nextMonth), Balance: 2.2, Change: 0},
				{Date: monthYear(now), Balance: 1.2, Change: 1.0},
			},
		},
		{ // Middle of the month, no balance.
			ops:  Operations{NewOperation(Deposit, 1.0, now)},
			summ: []Summary{{Date: monthYear(now), Balance: 0, Change: 1.0}},
		},
		{ // Two balances --> Use the most recent.
			ops: Operations{
				NewOperation(Balance, 1.2, incMin),
				NewOperation(Balance, 2.2, now),
			},
			summ: []Summary{{Date: monthYear(now), Balance: 1.2, Change: 0}},
		},
	}
	for _, d := range data {
		got := d.ops.Summarize()
		if (len(got) != 0 && len(d.summ) != 0) && !reflect.DeepEqual(got, d.summ) {
			t.Errorf("got:%+v want:%+v", got, d.summ)
		}
	}
}
