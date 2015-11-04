package isumm

import (
	"fmt"
	"net/http"
	"text/template"

	"appengine"
	"appengine/user"
)

var appTemplate = template.Must(template.ParseFiles("static/app.template.html"))

type appMonthlySummary struct {
	Header           string
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

	summ := make(map[MonthKey]*appMonthlySummary)
	var appInvestments []appInvestment
	for _, inv := range investments {
		appInvestments = append(appInvestments, appInvestment{Key: inv.Key, Name: inv.Name})
		for mKey, Sum := range inv.Ops.Summarize() {
			s, ok := summ[mKey]
			if !ok {
				s = &appMonthlySummary{
					Header:           fmt.Sprintf("%v/%v", mKey.Month, mKey.Year),
					MonthlySummaries: make([]*appInvestmentSummary, 0),
				}
				summ[mKey] = s
			}
			s.MonthlySummaries = append(s.MonthlySummaries, &appInvestmentSummary{Investment: inv.Name, Summary: Sum})
		}
	}

	var appMonthlySummaries []*appMonthlySummary
	for _, v := range summ {
		appMonthlySummaries = append(appMonthlySummaries, v)
	}
	if err := appTemplate.Execute(w, appParams{u.String(), logoutUrl, appInvestments, appMonthlySummaries}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
