package isumm

import (
	"net/http"
	"net/url"
	"time"

	"appengine"
	"appengine/urlfetch"
)

const (
	registryUrl = "http://isumm-registry.appspot.com//register"
	version     = "0.1"
)

var (
	throttleDuration = time.Second * 30
	// Every time a new instance is pushed, we would like to register it. Them start throttling.
	lastRegisterSent = time.Now().Add(-throttleDuration)
)

func SendRegisterRequest(c appengine.Context) {
	// If throttleDuration hasn't passed, no-op.
	if lastRegisterSent.After(time.Now().Add(-throttleDuration)) {
		return
	}
	lastRegisterSent = time.Now()
	// TODO(danielfireman): Parse response and use the returned version value to
	// print update request banner.
	client := &http.Client{
		Transport: &urlfetch.Transport{
			Context:                       c,
			Deadline:                      time.Second * 2,
			AllowInvalidServerCertificate: true,
		},
	}
	_, err := client.PostForm(registryUrl, url.Values{"version": {version}})
	if err != nil {
		c.Errorf("Failed sending register %q request.", err)
		return
	}
}
