package isumm

import (
	"fmt"
	"sort"
	"time"
)

// TimeseriesChart represents a chart where the x-axis is a timeseries.
type TimeseriesChart []timeSeriesPoint

// timeSeriesPoint represents a chart point where x-axis is a timeseries.
// All fields are private because interesting result is the String().
type timeSeriesPoint struct {
	date  time.Time
	value float32
}

func (g timeSeriesPoint) String() string {
	return fmt.Sprintf("[%d, %.2f]", g.date.UnixNano()/1000000, g.value)
}

func AmountSummaryChart(invs []*Investment) TimeseriesChart {
	summs := AggregateByDate(invs)
	var chart TimeseriesChart
	for _, s := range summs {
		chart = append(chart, timeSeriesPoint{s.Date, s.Balance})
	}
	return chart
}

func InterestRateChart(invs []*Investment) TimeseriesChart {
	summs := AggregateByDate(invs)
	// Interest rates implies at least two reference points.
	if len(summs) < 2 {
		return TimeseriesChart{}
	}
	var chart TimeseriesChart
	for i, s := range summs {
		t := timeSeriesPoint{date: s.Date}
		switch {
		case i == 0:
			t.value = 0
		default:
			t.value = float32(s.Balance-s.Change-summs[i-1].Balance) / float32(s.Change+summs[i-1].Balance) * 100.0
		}
		chart = append(chart, t)
	}
	return chart
}

// Returns the aggregated summaries sorted by date.
func AggregateByDate(invs []*Investment) Summaries {
	summByDate := make(map[time.Time]Summary)
	for _, i := range invs {
		for _, s := range i.Ops.Summarize() {
			v := summByDate[s.Date]
			v.Date = s.Date
			v.Balance += s.Balance
			v.Change += s.Change
			summByDate[s.Date] = v
		}
	}
	var summs Summaries
	for _, s := range summByDate {
		summs = append(summs, s)
	}
	// More information at: https://github.com/danielfireman/isumm/issues/2
	sort.Sort(summs)
	switch len(summs) {
	case 1:
		if summs[0].Balance == 0 {
			// There only one point and it is zero. This probably means that it's the first month the user
			// and the month hasn't changed.
			return Summaries{}
		}
	default:
		// This prevents the corner case where the last month hasn't ended and there already deposits or
		// widrawals. There is a valid case for that where we have zeroed the balance because of a Widrawal.
		last := summs[len(summs)-1]
		previous := summs[len(summs)-2]
		if last.Balance == 0 && (previous.Balance+last.Change != 0) {
			summs = summs[:len(summs)-2]
		}
	}
	return summs
}
