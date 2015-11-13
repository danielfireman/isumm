package isumm

import (
	"fmt"
	"log"
	"sort"
	"time"
)

// TimeseriesChart represents a chart where the x-axis is a timeseries.
type TimeseriesChart []timeSeriesPoint

func (t TimeseriesChart) Len() int {
	return len(t)
}

func (t TimeseriesChart) Less(i, j int) bool {
	return t[i].date.Before(t[j].date)
}

func (t TimeseriesChart) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

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
	balanceByDate := make(map[time.Time]float32)
	for _, i := range invs {
		for _, s := range i.Ops.Summarize() {
			balanceByDate[s.Date] += s.Balance
		}
	}
	var chart TimeseriesChart
	for t, b := range balanceByDate {
		chart = append(chart, timeSeriesPoint{t, b})
	}
	// Very important because go map deliberate changes the map interation from
	// time to time.
	sort.Sort(chart)
	return chart
}

func InterestRateChart(invs []*Investment) TimeseriesChart {
	summByDate := make(map[time.Time]Summary)
	for _, i := range invs {
		log.Printf("Inv: %v", i)
		for _, s := range i.Ops.Summarize() {
			log.Printf("Summ: %v", s)
			v := summByDate[s.Date]
			v.Date = s.Date
			v.Balance += s.Balance
			v.Change += s.Change
			summByDate[s.Date] = v
		}
	}
	log.Printf("byDate: %v", summByDate)
	var summs Summaries
	for _, s := range summByDate {
		summs = append(summs, s)
	}
	// The formula needs sorted summs. More information at: https://github.com/danielfireman/isumm/issues/2
	sort.Sort(summs)
	var chart TimeseriesChart
	for i, s := range summs {
		t := timeSeriesPoint{date: s.Date}
		switch {
		case i == 0:
			t.value = 0
		default:
			t.value = float32(s.Balance-s.Change-summs[i-1].Balance) / float32(s.Change+summs[i-1].Balance)
		}
		chart = append(chart, t)
	}
	log.Printf("%v", chart)
	return chart
}
