package isumm

import (
	"reflect"
	"testing"
	"time"
)

var (
	now       = time.Now()
	nextMonth = now.Add(32 * 24 * time.Hour)
	incMin    = 10 * time.Minute
)

func TestSummarize(t *testing.T) {
	data := []struct {
		ops  Operations
		summ MonthlySummary
	}{
		{ // Simple case: One entry containing only balance.
			ops:  Operations{{Date: now, Value: 1.2, Type: Balance}},
			summ: MonthlySummary{MonthKey{now.Month(), now.Year()}: {Balance: 1.2, Change: 0}},
		},
		{ // Empty case.
			ops:  Operations{},
			summ: MonthlySummary{},
		},
		{
			ops: Operations{
				{Date: now.Add(incMin), Value: 1.2, Type: Balance},
				{Date: now, Value: 1.0, Type: Deposit},
				{Date: nextMonth, Value: 2.2, Type: Balance},
			},
			summ: MonthlySummary{
				MonthKey{nextMonth.Month(), nextMonth.Year()}: {Balance: 2.2, Change: 0},
				MonthKey{now.Month(), now.Year()}:             {Balance: 1.2, Change: 1.0},
			},
		},
		{ // Middle of the month, no balance.
			ops:  Operations{{Date: now, Value: 1.0, Type: Deposit}},
			summ: MonthlySummary{MonthKey{now.Month(), now.Year()}: {Balance: 0, Change: 1.0}},
		},
		{ // Two balances --> Use the most recent.
			ops: Operations{
				{Date: now.Add(incMin), Value: 1.2, Type: Balance},
				{Date: now, Value: 2.2, Type: Balance},
			},
			summ: MonthlySummary{MonthKey{now.Month(), now.Year()}: {Balance: 1.2, Change: 0}},
		},
	}
	for _, d := range data {
		got := Summarize(d.ops)
		if (len(got) != 0 && len(d.summ) != 0) && !reflect.DeepEqual(got, d.summ) {
			t.Errorf("got:%+v want:%+v", got, d.summ)
		}
	}
}
