package isumm

import (
	"fmt"
	"net/http"
	"sort"
	"text/template"
	"time"

	"appengine"
	"appengine/user"
)

var appTemplate = template.Must(template.ParseFiles("static/app.template.html"))

// All information needed to render investment information.
type appInvestments []appInvestment

func (m appInvestments) Len() int      { return len(m) }
func (m appInvestments) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

// Sort investment lexicographically for rendering the summaries and ops
// section.
func (m appInvestments) Less(i, j int) bool {
	return m[i].Investment.Name < m[j].Investment.Name
}

type appInvestment struct {
	Investment *Investment
	Summary    MonthlySummary
}

type appParams struct {
	User         string
	Currency     string
	LogoutURL    string
	Investments  appInvestments
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

	w.Header().Set("Content-Type", "text/html")
	investments, err := GetInvestments(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var appInv appInvestments
	for _, inv := range investments {
		appInv = append(appInv, appInvestment{inv, Summarize(inv.Ops)})
	}
	sort.Sort(appInv)
	params := appParams{
		User:         u.String(),
		Currency:     Currency,
		LogoutURL:    logoutUrl,
		Investments:  appInv,
		SummaryGraph: AmountSummaryChart(appInv)}
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

func AmountSummaryChart(ais appInvestments) []timeSeriesPoint {
	auxChart := make(map[time.Time]float32)
	for _, ai := range ais {
		for _, s := range ai.Summary {
			auxChart[s.Date] += s.Balance
		}
	}
	var chart []timeSeriesPoint
	for t, b := range auxChart {
		chart = append(chart, timeSeriesPoint{t, b})
	}
	return chart
}
