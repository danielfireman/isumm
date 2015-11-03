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
