package isumm

import (
	"net/http"

	"appengine"
	"appengine/user"
)

func LogoutURL(c appengine.Context, r *http.Request) (string, error) {
	url, err := user.LogoutURL(c, r.URL.String())
	if err != nil {
		return "", err
	}
	return url, nil
}
