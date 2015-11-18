package isumm

import (
	"net/http"
	"text/template"

	"appengine"
	"appengine/user"
)

var appTemplate = template.Must(template.ParseFiles("static/app.template.html"))

type appParams struct {
	User               string
	Currency           string
	LogoutURL          string
	Investments        []*Investment
	AmountSummaryChart TimeseriesChart
	InterestRateChart  TimeseriesChart
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

	// Send fire-and-forget register request.
	go SendRegisterRequest(c)

	investments, err := GetInvestments(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params := appParams{
		User:               u.String(),
		Currency:           Currency,
		LogoutURL:          logoutUrl,
		Investments:        investments,
		AmountSummaryChart: AmountSummaryChart(investments),
		InterestRateChart:  InterestRateChart(investments)}
	if err := appTemplate.Execute(w, params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
