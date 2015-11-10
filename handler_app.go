package isumm

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/user"
)

var appTemplate = template.Must(template.ParseFiles("static/app.template.html"))

type appParams struct {
	User         string
	Currency     string
	LogoutURL    string
	Investments  []*Investment
	SummaryGraph []timeSeriesPoint
}

func App(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c) // Login is mandatory on this page. No need to check nil value here.
	if !IsUserAllowed(u) {
		InvalidUserPage(c, w, r, u)
		return
	}
	logoutUrl, err := LogoutURL(c, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	investments, err := GetInvestments(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params := appParams{
		User:         u.String(),
		Currency:     Currency,
		LogoutURL:    logoutUrl,
		Investments:  investments,
		SummaryGraph: AmountSummaryChart(investments)}
	if err := appTemplate.Execute(w, params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// timeSeriesPoint represents a chart point where x-axis is a timeseries.
// All fields are private because interesting result is the String().
type timeSeriesPoint struct {
	date    time.Time
	balance float32
}

func (g timeSeriesPoint) String() string {
	return fmt.Sprintf("[%d, %.2f]", g.date.UnixNano()/1000000, g.balance)
}

func AmountSummaryChart(invs []*Investment) []timeSeriesPoint {
	auxChart := make(map[time.Time]float32)
	for _, i := range invs {
		for _, s := range i.Ops.Summarize() {
			auxChart[s.Date] += s.Balance
		}
	}
	var chart []timeSeriesPoint
	for t, b := range auxChart {
		chart = append(chart, timeSeriesPoint{t, b})
	}
	return chart
}
