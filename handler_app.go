package isumm

import (
	"fmt"
	"net/http"
	"sort"
	"text/template"

	"appengine"
	"appengine/user"
)

var appTemplate = template.Must(template.ParseFiles("static/app.template.html"))

type monthlySummaries []*appMonthlySummary

func (m monthlySummaries) Len() int      { return len(m) }
func (m monthlySummaries) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m monthlySummaries) Less(i, j int) bool {
	if m[i].MonthKey.Year != m[j].MonthKey.Year {
		return m[i].MonthKey.Year < m[j].MonthKey.Year
	}
	return m[i].MonthKey.Month < m[j].MonthKey.Month
}

type appMonthlySummary struct {
	Header           string
	MonthKey         MonthKey
	MonthlySummaries []*appInvestmentSummary
}

type appInvestmentSummary struct {
	Investment string
	Summary    Summary
}

type appInvestment struct {
	Key  string
	Name string
}

type appParams struct {
	User         string
	LogoutUrl    string
	Investments  []appInvestment
	AllSummaries []*appMonthlySummary
	SummaryGraph []graphPoint
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

	// This is a way to present the summary in a month-by-month basis.
	summ := make(map[MonthKey]*appMonthlySummary)
	var appInvestments []appInvestment
	for _, inv := range investments {
		appInvestments = append(appInvestments, appInvestment{Key: inv.Key, Name: inv.Name})
		for mKey, Sum := range inv.Ops.Summarize() {
			s, ok := summ[mKey]
			if !ok {
				s = &appMonthlySummary{
					MonthKey:         mKey,
					Header:           fmt.Sprintf("%v/%v", mKey.Month, mKey.Year),
					MonthlySummaries: make([]*appInvestmentSummary, 0),
				}
				summ[mKey] = s
			}
			s.MonthlySummaries = append(s.MonthlySummaries, &appInvestmentSummary{Investment: inv.Name, Summary: Sum})
		}
	}

	var appMonthlySummaries monthlySummaries
	for _, v := range summ {
		appMonthlySummaries = append(appMonthlySummaries, v)
	}
	sort.Sort(appMonthlySummaries)
	if err := appTemplate.Execute(w, appParams{u.String(), logoutUrl, appInvestments, appMonthlySummaries, GraphSummary(appMonthlySummaries)}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type graphPoint struct {
	x int
	y float32
}

func (g graphPoint) String() string {
	return fmt.Sprintf("[%d, %.2f]", g.x, g.y)
}

func GraphSummary(summary monthlySummaries) []graphPoint {
	var graph []graphPoint
	for monthIndex, ms := range summary {
		summ := float32(0)
		for _, is := range ms.MonthlySummaries {
			summ += is.Summary.Balance
		}
		graph = append(graph, graphPoint{monthIndex, summ})
	}
	return graph
}
