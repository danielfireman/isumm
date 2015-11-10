package isumm

import (
	"sort"
	"time"
)

const (
	opEntityName = "Operation"
)

type OpType int

const (
	Unknown = iota
	Balance
	Deposit
	Withdrawal
)

func (t OpType) String() string {
	switch t {
	case 1:
		return "Balance"
	case 2:
		return "Deposit"
	case 3:
		return "Withdrawal"
	}
	return "Unknown"
}

type Operation struct {
	Value float32   `datastore:"value,noindex"`
	Date  time.Time `datastore:"date`
	Type  OpType    `datastore:"type"`
}

type Operations []Operation

func (ops Operations) Len() int {
	return len(ops)
}

func (ops Operations) Less(i, j int) bool {
	return ops[i].Date.After(ops[j].Date)
}

func (ops Operations) Swap(i, j int) {
	aux := ops[i]
	ops[i] = ops[j]
	ops[j] = aux
}

type Summary struct {
	Date    time.Time
	Balance float32
	Change  float32
}

func (s *Summary) Reset() {
	s.Balance = 0
	s.Change = 0
}

func (ops Operations) Summarize() []Summary {
	if !sort.IsSorted(ops) {
		sort.Sort(ops)
	}
	var returned []Summary
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
