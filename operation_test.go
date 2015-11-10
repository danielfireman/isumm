package isumm

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

var (
	now       = time.Now()
	nextMonth = now.Add(32 * 24 * time.Hour)
	incMin    = 10 * time.Minute
)

func TestSort(t *testing.T) {
	o := Operations{{Date: now}, {Date: now.Add(incMin)}}
	sort.Sort(o)
	if o[0].Date.Before(o[1].Date) {
		t.Errorf("want %s before than %s", o[0], o[1])
	}
}

func TestSummarize(t *testing.T) {
	data := []struct {
		ops  Operations
		summ []Summary
	}{
		{ // Simple case: One entry containing only balance.
			ops:  Operations{{Date: now, Value: 1.2, Type: Balance}},
			summ: []Summary{{Date: monthYear(now), Balance: 1.2, Change: 0}},
		},
		{ // Empty case.
			ops:  Operations{},
			summ: []Summary{},
		},
		{
			ops: Operations{
				{Date: now.Add(incMin), Value: 1.2, Type: Balance},
				{Date: now, Value: 1.0, Type: Deposit},
				{Date: nextMonth, Value: 2.2, Type: Balance},
			},
			summ: []Summary{
				{Date: monthYear(nextMonth), Balance: 2.2, Change: 0},
				{Date: monthYear(now), Balance: 1.2, Change: 1.0},
			},
		},
		{ // Middle of the month, no balance.
			ops:  Operations{{Date: now, Value: 1.0, Type: Deposit}},
			summ: []Summary{{Date: monthYear(now), Balance: 0, Change: 1.0}},
		},
		{ // Two balances --> Use the most recent.
			ops: Operations{
				{Date: now.Add(incMin), Value: 1.2, Type: Balance},
				{Date: now, Value: 2.2, Type: Balance},
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
