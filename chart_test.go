package isumm

import (
	"reflect"
	"testing"
	"time"
)

func TestAggregateByDate(t *testing.T) {
	timestamp := int64(1)
	monthYearTimestamp := monthYear(time.Unix(timestamp, 0))
	testData := []struct {
		desc string
		invs []*Investment
		want Summaries
	}{
		{
			"no investments",
			[]*Investment{},
			Summaries{},
		},
		{
			"first month-noops",
			[]*Investment{{Key: "1", Name: "Foo"}},
			Summaries{},
		},
		{
			"one inv",
			[]*Investment{
				{Key: "1", Name: "Foo", Ops: []Operation{
					{Value: 10, Type: Deposit, Timestamp: timestamp},
					{Value: 11, Type: Balance, Timestamp: timestamp + 1},
				}},
			},
			Summaries{{Date: monthYearTimestamp, Balance: 11, Change: 10}},
		},
	}
	for _, test := range testData {
		got := AggregateByDate(test.invs)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("(%s) got:%v want:%v invs:%v", test.desc, got, test.want, test.invs)
		}
	}
}
