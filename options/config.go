package options

import (
	"time"
)

type Config struct {
	LogLevel        string
	Port            int
	HealthCheckPort int
	PrometheusPort  int
	KeyID           string
	AuthPageLink    string
	PrivateKeyBytes []byte
	IDPConfFileName string
	TTL             time.Duration
}
