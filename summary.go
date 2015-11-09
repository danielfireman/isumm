package isumm

import (
	"sort"
	"time"
)

type Summary struct {
	Balance float32
	Change  float32
}

func (s *Summary) Reset() {
	s.Balance = 0
	s.Change = 0
}

type MonthKey struct {
	Month time.Month
	Year  int
}

type MonthlySummary map[MonthKey]Summary

func Summarize(ops Operations) MonthlySummary {
	if !sort.IsSorted(ops) {
		sort.Sort(ops)
	}
	returned := make(MonthlySummary)
	var summ Summary
	var month time.Month
	var year int
	for i, op := range ops {
		if i > 0 && ops[i-1].Date.Month() != op.Date.Month() {
			returned[MonthKey{month, year}] = summ
			summ.Reset()
		}
		month = op.Date.Month()
		year = op.Date.Year()
		switch op.Type {
		case Withdrawal:
			summ.Change -= op.Value
		case Deposit:
			summ.Change += op.Value
		case Balance:
			// This enforces the case when we have more than one balance, the
			// system will use the most recent. This is easier than change
			// loads of places to make sure enforce only one balance.
			if summ.Balance == 0 {
				summ.Balance = op.Value
			}
		}
	}
	// We must not forget the latest summary.
	if len(ops) > 0 {
		returned[MonthKey{month, year}] = summ
	}
	return returned
}
