package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/akhripko/auth2-jwt/metrics"
	"github.com/akhripko/auth2-jwt/options"
	"github.com/akhripko/auth2-jwt/srv/healthcheck"
	"github.com/akhripko/auth2-jwt/srv/prometheus"
	"github.com/akhripko/auth2-jwt/srv/srvhttp"
	localstorage "github.com/akhripko/auth2-jwt/storage/local"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	// read service config from os env
	config := options.ReadEnv()

	// init logger
	initLogger(config)

	log.Info("begin...")
	// register metrics
	metrics.Register()

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	var wg = &sync.WaitGroup{}

	// load private key
	keyData, _ := pem.Decode(config.PrivateKeyBytes)
	if keyData == nil {
		log.Error("jwt private key: pem decode failed")
		os.Exit(1)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyData.Bytes)
	if err != nil {
		log.Error("failed to parse private key: " + err.Error())
		os.Exit(1)
	}

	// build ipd configs storage
	ipdConfigsStorage, err := buildIDPConfigsStorage(ctx, config)
	if err != nil {
		log.Error("ipd configs storage init error:", err.Error())
		os.Exit(1)
	}

	loginPage, err := readFile("login-page.html")
	if err != nil {
		log.Error("failed to load login page:", err.Error())
		os.Exit(1)
	}

	// build main http srv
	httpSrv, err := srvhttp.New(config.Port, loginPage, ipdConfigsStorage, privateKey, config.KeyID, config.AuthPageLink, config.TTL)
	if err != nil {
		log.Error("http srv init error:", err.Error())
		os.Exit(1)
	}

	// build prometheus srv
	prometheusSrv := prometheus.New(config.PrometheusPort)
	// build healthcheck srv
	healthSrv := healthcheck.New(
		config.HealthCheckPort,
		prometheusSrv.Check,
		httpSrv.Check,
		ipdConfigsStorage.Check,
	)

	// run srv
	healthSrv.Run(ctx, wg)
	prometheusSrv.Run(ctx, wg)
	httpSrv.Run(ctx, wg)

	// wait while services work
	wg.Wait()
	log.Info("end")
}

func initLogger(config options.Config) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	switch strings.ToLower(config.LogLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Error("Got Interrupt signal")
		stop()
	}()
}

func buildIDPConfigsStorage(ctx context.Context, config options.Config) (srvhttp.IDPConfigsStorage, error) {
	// try use back-office
	//if len(config.BackofficeConfig.Target) > 0 {
	//	boClient, err := backoffice.New(ctx, config.BackofficeConfig)
	//	if err != nil {
	//		return nil, errors.Wrap(err, "failed to init domain storage")
	//	}
	//	return boClient, nil
	//}

	// try use config from yaml file
	idpConfBytes, err := readFile(config.IDPConfFileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config from yaml file")
	}
	inmemStorage, err := localstorage.BuildInMemStorageFromYaml(idpConfBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build storage from yaml")
	}
	return inmemStorage, nil
}

func readFile(fileName string) ([]byte, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file '%s'", fileName)
	}
	return data, nil
}
