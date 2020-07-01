package local

import (
	"testing"

	"github.com/akhripko/auth2-jwt/models"
	"github.com/stretchr/testify/assert"
)

func TestBuildFromYaml(t *testing.T) {
	yamlBytes := []byte(`
host:
  ad:
    idptype: microsoft
    servicetype: staff
    account: blablaaccount
    audience: blablaaudience
    auth:
      clientid: blablaclientid
      clientsecret: blablaclientsecret
      endpoint:
        authurl: https://login.microsoftonline.com/common/oauth2/v2.0/authorize
        tokenurl: https://login.microsoftonline.com/common/oauth2/v2.0/token
        authstyle: 0
      redirecturl: srvhttp://localhost:8080/ad/callback
      scopes:
      - openid
    endpointprofile: https://graph.microsoft.com/v1.0/me`)

	store, err := BuildInMemStorageFromYaml(yamlBytes)
	assert.NoError(t, err)
	accConf, err := store.ReadAccountConfig("host")
	assert.NoError(t, err)
	authConf, err := accConf.GetAuthConfig("ad")
	assert.NoError(t, err)
	assert.Equal(t, models.Microsoft, authConf.IDPType)
	assert.Equal(t, models.StaffServiceType, authConf.ServiceType)
	assert.Equal(t, "https://graph.microsoft.com/v1.0/me", authConf.EndPointProfile)
	assert.Equal(t, "blablaclientid", authConf.Auth.ClientID)
	assert.Equal(t, "blablaaudience", authConf.Audience)
	assert.Equal(t, "blablaaccount", authConf.Account)
}

func TestCheck(t *testing.T) {
	store := NewInMemStorage(nil)
	assert.NotNil(t, store)
	assert.Error(t, store.Check())
}
