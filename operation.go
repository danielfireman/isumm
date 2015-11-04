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

func (ops Operations) FilterMonth(month time.Month, year int) Operations {
	var monthOps Operations
	for _, o := range ops {
		if o.Date.Month() == month && o.Date.Year() == year {
			monthOps = append(monthOps, o)
		}
	}
	return monthOps
}

func (ops Operations) Summarize() MonthlySummary {
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
			summ.Balance = op.Value
		}
	}
	// We must not forget the latest summary.
	if len(ops) > 0 {
		returned[MonthKey{month, year}] = summ
	}
	return returned
}
