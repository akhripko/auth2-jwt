package srvhttp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/akhripko/auth2-jwt/jwt"
	"github.com/akhripko/auth2-jwt/models"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type (
	link struct {
		Title string `json:"title"`
		Link  string `json:"link"`
	}

	linksListResponse struct {
		Links            []link `json:"buttons"`
		ManningAgentLink string `json:"manningAgentLink,omitempty"`
	}
)

// /?state=...&redirect_uri=...
func (s *Service) handleMain(w http.ResponseWriter, r *http.Request) {
	account := getAccount(r)
	accConf, err := getAccConf(account, s.storage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("main: read account config err:", err)
		return
	}
	if len(accConf) == 1 {
		conf, err := accConf.GetOneAuthConfig()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("auth: read auth config err:", err)
			return
		}
		makeIDPRedirect(w, r, conf)
		return
	}

	http.Redirect(w, r,
		s.AuthPageLink+"?"+r.URL.Query().Encode(),
		http.StatusTemporaryRedirect)
}

// ?state=...&redirect_uri=...
func (s *Service) handleAuthHTMLPage(w http.ResponseWriter, r *http.Request) {
	w.Write(s.loginPage)
}

// /links
// /links?state=...&redirect_uri=...
// support CORS requests
func (s *Service) handleGetAuthLinks(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get account config
	account := getAccount(r)
	accConf, err := getAccConf(account, s.storage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("get auth links: read auth config err:", err)
		return
	}

	res := &linksListResponse{}
	for idp, auth := range accConf {
		if auth.ServiceType == models.ManningAgentServiceType {
			res.ManningAgentLink = buildIDPAuthEndpoint(idp, r)
			continue
		}

		res.Links = append(res.Links, link{
			Title: auth.Title,
			Link:  buildIDPAuthEndpoint(idp, r),
		})
	}

	// write json linksListResponse
	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("get auth links: build json linksListResponse err:", err)
		return
	}
	writeJSONBytes(w, data)
}

// /{IDP}/authorize?state=...&redirect_uri=...
func (s *Service) handleAuth(w http.ResponseWriter, r *http.Request) {
	account := getAccount(r)
	idp := mux.Vars(r)["IDP"]
	authConf, err := getAuthConf(account, idp, s.storage)
	if err != nil {
		rootRedirect(w, r)
		log.Error("authorize: read auth config err:", err)
		return
	}
	makeIDPRedirect(w, r, authConf)
}

// /{IDP}/callback?state=...&code=...
// state is base64 string: redirect_uri|state|account
func (s *Service) handleCallback(w http.ResponseWriter, r *http.Request) {
	const userTokenType = "user"

	// restore client state
	clState, err := getClientState(r)
	if err != nil {
		rootRedirect(w, r)
		log.Error("idp callback: get client state err:", err)
		return
	}
	account := getAccount(r)
	idp := mux.Vars(r)["IDP"]
	authConf, err := getAuthConf(account, idp, s.storage)
	if err != nil {
		rootRedirect(w, r)
		log.Error("idp callback: read auth config err:", err)
		return
	}

	// exchange code to token
	q := r.URL.Query()
	code := q.Get("code")
	token, err := exchangeCodeToToken(authConf, code)
	if err != nil {
		rootRedirect(w, r)
		log.Error("idp callback: idp exchange err:", err.Error())
		return
	}

	// build user model
	user, err := buildUserModel(authConf, token.AccessToken)
	if err != nil {
		rootRedirect(w, r)
		log.Error("idp callback err: ", err.Error())
		return
	}

	// build token claims
	t := time.Now().UTC()
	ttl := s.ttl
	if authConf.TTL != nil {
		ttl = *authConf.TTL
	}

	claims := jwt.Claims{
		UserID:     user.ID,
		Account:    authConf.Account,
		Type:       userTokenType,
		Name:       user.Name,
		Email:      user.Email,
		Nickname:   user.NickName,
		GivenName:  user.GivenName,
		FamilyName: user.FamilyName,
		Picture:    user.Picture,
		Groups:     user.Groups,
		Context:    user.Context,
		StandardClaims: jwtgo.StandardClaims{
			Issuer:    r.Host,
			Audience:  authConf.Audience,
			IssuedAt:  t.Unix(),
			ExpiresAt: t.Add(ttl).Unix(),
		},
	}

	// build jwt token
	tokenStr, err := jwt.BuildSignedToken(s.rsaPrivateKey, s.keyID, claims)
	if err != nil {
		rootRedirect(w, r)
		log.Error("idp callback: build token err: ", err.Error())
		return
	}

	callbackRedirect(w, r, &clState, tokenStr)
}

// /key
func (s *Service) getPublicKey(w http.ResponseWriter, r *http.Request) {
	res, _ := jwt.ExportRsaPublicKeyAsPemStr(&s.rsaPrivateKey.PublicKey)
	w.Write([]byte(res))
}
