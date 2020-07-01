package srvhttp

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/akhripko/auth2-jwt/idp"
	"github.com/akhripko/auth2-jwt/models"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func toBase64(src []byte) []byte {
	data := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(data, src)
	return data
}

func fromBase64(src []byte) ([]byte, error) {
	data := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(data, src)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func makeIDPRedirect(w http.ResponseWriter, r *http.Request, authConf *models.AuthConfig) {
	http.Redirect(w, r, authConf.Auth.AuthCodeURL(buildClientStateString(r)), http.StatusTemporaryRedirect)
}

func buildIDPAuthEndpoint(idp string, r *http.Request) string {
	q := r.URL.Query()
	return fmt.Sprintf("//%s/auth2/%s/authorize?%s", r.Host, idp, q.Encode())
}

func rootRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth2", http.StatusTemporaryRedirect)
}

func callbackRedirect(w http.ResponseWriter, r *http.Request, clState *ClientState, tokenStr string) {
	url := clState.GetRedirectURI()
	if strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "#id_token=" + tokenStr + "&state=" + clState.GetInitialState()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func buildUserModel(authConf *models.AuthConfig, accessToken string) (*models.User, error) {
	// read user info
	req, err := http.NewRequest("GET", authConf.EndPointProfile, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read user info")
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read user info")
	}
	defer res.Body.Close()

	userInfo, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read user info")
	}

	// build user model
	buildUser, err := idp.GetUserBuilder(authConf.IDPType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get builder")
	}
	user, err := buildUser(userInfo, accessToken, authConf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build user model")
	}

	return user, nil
}

func exchangeCodeToToken(authConf *models.AuthConfig, code string) (*oauth2.Token, error) {
	// exchange code to token
	token, err := authConf.Auth.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange idp token")
	}
	if !token.Valid() {
		return nil, errors.Wrap(err, "idp token is not valid")
	}
	return token, nil
}

func getAccount(r *http.Request) string {
	domains := strings.Split(r.Host, ".")
	return domains[0]
}

func getAccConf(account string, storage IDPConfigsStorage) (models.AccountConfig, error) {
	// get account config
	accConf, err := storage.ReadAccountConfig(account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read account config")
	}

	return accConf, nil
}

func getAuthConf(account, idp string, storage IDPConfigsStorage) (*models.AuthConfig, error) {
	// get account config
	accConf, err := getAccConf(account, storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read account config")
	}
	// get idp auth config
	authConf, err := accConf.GetAuthConfig(idp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read auth config")
	}

	return authConf, nil
}

func writeJSONBytes(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
