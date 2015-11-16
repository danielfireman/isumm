package isumm

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
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

func (t OpType) IsValid() bool {
	return t.String() != "Unknown"
}

type Operation struct {
	Value     float32 `datastore:"value,noindex"`
	Type      OpType  `datastore:"type"`
	Timestamp int64   `datastore:"date"`
}

// TODO(danielfireman): Verify if the date needs to be cached.
func (o *Operation) Date() time.Time {
	return time.Unix(0, o.Timestamp)
}

func (o *Operation) SetDate(d time.Time) {
	o.Timestamp = d.UnixNano()
}

func NewOperation(t OpType, v float32, d time.Time) Operation {
	return Operation{Type: t, Value: v, Timestamp: d.UnixNano()}
}

func NewOperationFromString(t, v, d string) (Operation, error) {
	opTypeInt, err := strconv.Atoi(t)
	if err != nil {
		return Operation{}, fmt.Errorf("Invalid operation type: \"%s\"", t)
	}
	opType := OpType(opTypeInt)
	if !opType.IsValid() {
		return Operation{}, fmt.Errorf("Invalid operation type: \"%s\"", t)
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(v), 32)
	if err != nil {
		return Operation{}, fmt.Errorf("Invalid value: \"%s\"", v)
	}
	date, err := time.Parse("2006-01-02", d)
	if err != nil {
		return Operation{}, fmt.Errorf("Invalid operation date: %s", d)
	}
	return Operation{Timestamp: date.UnixNano(), Value: float32(value), Type: opType}, nil
}

type Operations []Operation

func (ops Operations) Len() int {
	return len(ops)
}

func (ops Operations) Less(i, j int) bool {
	return ops[i].Timestamp > ops[j].Timestamp
}

func (ops Operations) Swap(i, j int) {
	ops[i], ops[j] = ops[j], ops[i]
}

type SummaryOp struct {
	Index     int
	Operation Operation
}

type Summary struct {
	Date       time.Time
	Balance    float32
	Change     float32
	SummaryOps []SummaryOp
}

func (s *Summary) Reset() {
	s.Balance = 0
	s.Change = 0
	s.SummaryOps = nil
}

type Summaries []Summary

func (s Summaries) Len() int {
	return len(s)
}

func (s Summaries) Less(i, j int) bool {
	return s[i].Date.Before(s[j].Date)
}

func (s Summaries) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (ops Operations) Summarize() Summaries {
	if !sort.IsSorted(ops) {
		sort.Sort(ops)
	}
	var returned []Summary
	var summ Summary
	for i, op := range ops {
		// If the month has changed (please notice that ops are sorted by
		// month).
		if i > 0 && ops[i-1].Date().Month() != op.Date().Month() {
			returned = append(returned, summ)
			summ.Reset()
		}
		// At this point we are only interested on the month/year pair and
		// would like to enjoy the datetime goodies and format (flot and
		// other libraries support it natively). So, setting the day to
		// the first day of the month.
		summ.Date = monthYear(op.Date())
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
		// We aggregate the operations for that month as part of the Operations
		// field, which allows us to show a nice log of everything affecting that
		// month.
		summ.SummaryOps = append(summ.SummaryOps, SummaryOp{Index: i, Operation: op})
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
