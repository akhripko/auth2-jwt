package srvhttp

import (
	"net/http"
	"net/url"
	"strings"
)

type ClientState []string

func (c ClientState) GetRedirectURI() string {
	if c == nil {
		return ""
	}
	if len(c) == 0 {
		return ""
	}
	return c[0]
}

func (c ClientState) GetInitialState() string {
	if c == nil {
		return ""
	}
	if len(c) < 2 {
		return ""
	}
	return c[1]
}

const stateSep = "|"

func buildClientStateString(r *http.Request) string {
	q := r.URL.Query()
	return url.QueryEscape(string(toBase64([]byte(
		q.Get("redirect_uri") + stateSep + q.Get("state")))))
}

func getClientState(r *http.Request) (ClientState, error) {
	state, err := url.QueryUnescape(r.URL.Query().Get("state"))
	if err != nil {
		return nil, err
	}
	stateBytes, err := fromBase64([]byte(state))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(stateBytes), stateSep), nil
}
