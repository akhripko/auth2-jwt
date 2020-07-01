package srvhttp

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	state := []byte("abc|||123")
	eState := toBase64(state)
	state2, err := fromBase64(eState)
	assert.NoError(t, err)
	assert.Equal(t, state, state2)
}

func TestBuildAccAuthEndpoint(t *testing.T) {
	r, err := http.NewRequest("GET", "http://host:8080/auth2?state=123&redirect_uri=https%3A%2F%2Fmy.drydock.studio", nil)
	assert.NoError(t, err)
	res := buildIDPAuthEndpoint("acc-idp", r)
	assert.Equal(t, "//host:8080/auth2/acc-idp/authorize?redirect_uri=https%3A%2F%2Fmy.drydock.studio&state=123", res)
}

func Test_getAccount(t *testing.T) {
	res := getAccount(&http.Request{Host: "acc1.auth.com"})
	assert.Equal(t, "acc1", res)

	res = getAccount(&http.Request{Host: "acc1"})
	assert.Equal(t, "acc1", res)

	res = getAccount(&http.Request{})
	assert.Equal(t, "", res)
}
