package srvhttp

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClientState(t *testing.T) {
	r, err := http.NewRequest("GET", "http://blabla?state=123&redirect_uri=https%3A%2F%2Fmy.drydock.studio", nil)
	assert.NoError(t, err)
	clState := buildClientStateString(r)

	r2, err := http.NewRequest("GET", "http://blabla?state="+clState, nil)
	assert.NoError(t, err)

	state, err := getClientState(r2)
	assert.NoError(t, err)
	assert.Equal(t, "https://my.drydock.studio", state.GetRedirectURI())
	assert.Equal(t, "123", state.GetInitialState())
}
