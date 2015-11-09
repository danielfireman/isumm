package isumm

import "time"

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

func (ops Operations) FilterMonth(month time.Month, year int) Operations {
	var monthOps Operations
	for _, o := range ops {
		if o.Date.Month() == month && o.Date.Year() == year {
			monthOps = append(monthOps, o)
		}
	}
	return monthOps
}
