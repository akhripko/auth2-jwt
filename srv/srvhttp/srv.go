package srvhttp

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type providerIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func New(port int, loginPage []byte, storage IDPConfigsStorage, key *rsa.PrivateKey, keyID, authPageLink string, ttl time.Duration) (*Service, error) {
	// build http server
	httpSrv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	// build Service
	var srv Service
	srv.setupHTTP(&httpSrv)
	srv.storage = storage
	srv.rsaPrivateKey = key
	srv.keyID = keyID
	srv.AuthPageLink = authPageLink
	srv.ttl = ttl
	srv.loginPage = loginPage

	return &srv, nil
}

type Service struct {
	http          *http.Server
	runErr        error
	readiness     bool
	storage       IDPConfigsStorage
	rsaPrivateKey *rsa.PrivateKey
	keyID         string
	AuthPageLink  string
	providerIndex *providerIndex
	ttl           time.Duration
	loginPage     []byte
}

func (s *Service) setupHTTP(srv *http.Server) {
	srv.Handler = s.buildHandler()
	s.http = srv
}

func (s *Service) buildHandler() http.Handler {
	r := mux.NewRouter()
	// path -> handlers

	r.HandleFunc("/auth2", s.handleMain).Methods("GET")
	r.HandleFunc("/auth2/authorize", s.handleMain).Methods("GET")
	r.HandleFunc("/auth2/login.html", s.handleAuthHTMLPage).Methods("GET")
	r.HandleFunc("/auth2/links", s.handleGetAuthLinks).Methods("GET")
	r.HandleFunc("/auth2/{IDP}/authorize", s.handleAuth).Methods("GET")
	r.HandleFunc("/auth2/{IDP}/callback", s.handleCallback).Methods("GET")
	r.HandleFunc("/auth2/key", s.getPublicKey).Methods("GET")

	// ==============
	return r
}

func (s *Service) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("http srv: begin run")

	go func() {
		defer wg.Done()
		log.Debug("http srv addr:", s.http.Addr)
		err := s.http.ListenAndServe()
		s.runErr = err
		log.Info("http srv: end run (", err, ")")
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		err := s.http.Shutdown(sdCtx)
		if err != nil {
			log.Info("http srv shutdown (", err, ")")
		}
	}()

	s.readiness = true
}

func (s *Service) Check() error {
	if !s.readiness {
		return errors.New("http srv is't ready yet")
	}
	if s.runErr != nil {
		return errors.New("run http srv issue")
	}
	return nil
}
