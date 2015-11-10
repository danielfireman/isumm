package isumm

import (
	"sort"
	"time"
)

type Summary struct {
	Date    time.Time
	Balance float32
	Change  float32
}

func (s *Summary) Reset() {
	s.Balance = 0
	s.Change = 0
}

type MonthlySummary []Summary

func (m MonthlySummary) Len() int           { return len(m) }
func (m MonthlySummary) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MonthlySummary) Less(i, j int) bool { return m[i].Date.After(m[j].Date) }

func Summarize(ops Operations) MonthlySummary {
	if !sort.IsSorted(ops) {
		sort.Sort(ops)
	}
	var returned MonthlySummary
	var summ Summary
	for i, op := range ops {
		// If the month has changed (please notice that ops are sorted by
		// month).
		if i > 0 && ops[i-1].Date.Month() != op.Date.Month() {
			returned = append(returned, summ)
			summ.Reset()
		}
		// At this point we are only interested on the month/year pair and
		// would like to enjoy the datetime goodies and format (flot and
		// other libraries support it natively). So, setting the day to
		// the first day of the month.
		summ.Date = monthYear(op.Date)
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
		returned = append(returned, summ)
	}
	return returned
}

func monthYear(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 01, 0, 0, 0, 0, t.Location())
}
