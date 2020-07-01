package options

import (
	"github.com/spf13/viper"
)

func ReadEnv() Config {
	viper.AutomaticEnv()

	viper.SetEnvPrefix("APP")

	viper.SetDefault("LOG_LEVEL", "DEBUG")
	viper.SetDefault("HTTP_PORT", 8080)
	viper.SetDefault("HEALTH_CHECK_PORT", 8888)
	viper.SetDefault("PROMETHEUS_PORT", 9100)
	viper.SetDefault("KEY_ID", "auth2-service")
	viper.SetDefault("AUTH_PAGE_LINK", "auth2/login.html")
	viper.SetDefault("ID_TOKEN_TTL", "24h")

	//viper.SetDefault("BACKOFFICE_TARGET", "back-office-service")
	viper.SetDefault("BACKOFFICE_RUN_JOB_INTERVAL", "10m")

	viper.SetDefault("IDP_CONFIG_FILE_NAME", "idpconf.yml")

	// FOR TEST
	viper.SetDefault("JWT_PRIVATE_KEY", `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA1ZDOfs1Am0nXfHHU2qloFn4Ak7FsdvsADv3Mw+rMraIkLvy+
ewWQEKQztKPz82bVITotrg3M6zBTUtIGusQEHynAdL/LZTyqS8u3Ji3EV2/zEUXL
JKl1XH0KSkaXlunKDtSPgziqy/UZDfVFXSU+9elVO9KmLwr2b9YK+RHIOA/YvkzJ
eSr17GNJHa5194tw8OiGgInIm2HUmQ06ffJueMZuqcK3n5JTOC0UJTlu7xhJAg64
9vEJ0o9TCQzGmriJvXvYKu/9RUvPFyAmscLzqGPuKBQU/No63sW5DNbr766KvbN9
2oR9ZmuMlInIj80d+jE03YO5qBzu3m/Mgss/dQIDAQABAoIBABWa7qj5Tr3m/HvE
cEomtTcBTEGkw1pODjV2C48OY+J08sAmJYcAixCD7A8sPvOyxYD1B8YB82cfnm5U
fQXL6rbUXHnzJTp6pqfAVijwnFpUIvzuWJy+3/aZV78n4RINWLmjW1llFDJJ+1zk
IT5JFFOAZXUF3HYJ3+B0gOutltnE4UUw5mgjhQPKuegZgpxxqUAx/3OZghj0XpSi
BkgtLjfURZN3GLqSJk+2C3RGihGDok3jGsij2ZvWfh6XuTncDLg7mx1U2qUOLRcO
NaHRu5/nUoTbCj7dZHYiow6x7sFsm0AS6kSmDgc4GihY8OCx+xWCCE223kErigHH
HjbFgkECgYEA9VUWSe6QM6v8wgptO3ZPy4vDZhla7AbHwdH6jNIBBA0CbEAf5J5i
m5pVlC51J+dccSUXT0zfDIL2qwkA11BjBv48wrjSme/09XT3khepR0YwMEoONZSb
+Yjt/z7yI8awITKwN6jdvzyq/2XNZTPVD9UDM7FpURsnmN8v+PLjZo0CgYEA3tob
3wc/gwf1jzqpdFkQzFJojrIN0jbpA9WlI25+LmDbqcTP/kNpSvd8aCqVspop0XjU
LCBG4PtDUxaaMI+1tdviXLAgoeOONvNZSGFJkv2ecy8gPbPIb7A6/aIlcLB4eoGV
wZ2tIFgBNnEbU8YerRfuYPDTmDJhtsk5zZ6+VokCgYEAseYAYp4WfRGzGHX9HYg+
dqgjig7KyqpkAd0k5SaHTAuu/RZGyj473P++HNTPaZ8wfm9aBswEVgtFmWLO9FM7
a9/B6aWiObQyGie2On7j5mY0HWAmC552uC0d/+ACMAUsxPX/qGzQV5NDoC2PElrS
nkdlqf91EjYxsX3uSITAdyUCgYEA2vf+yM9rZYmNjDW8yVi2e71BhWyIzhQsMxta
zwzDqTh8vjxnjtOYAxRYRlYJj1uRWYTbHZx9aJUa2upriOm8RzwOGLrq5Ycddvr3
sHn/fBH9/fnBOT+M48mKvSr0lNyhFOZ9Sqhus0glsOPEUTVrcPMBxHj9wB9JCfyA
8nxXU9kCgYEA4ZRiZ7iFI4nQIHyau0oOIOLf27UljV5AHVAUs2+W4pSG+ETpC1W5
/vODz+sLWFPDhP9s+bU56qGlKslM6Etyqi/5uEg8gm17mNozaGhF75m4yaL6oPhv
ObL0fLzoBYfeRFkYeIYxeCY8y4J2YzOBwX9xNNfQCRBwjfHhvy6rrSU=
-----END RSA PRIVATE KEY-----`)

	return Config{
		LogLevel:        viper.GetString("LOG_LEVEL"),
		Port:            viper.GetInt("HTTP_PORT"),
		HealthCheckPort: viper.GetInt("HEALTH_CHECK_PORT"),
		PrometheusPort:  viper.GetInt("PROMETHEUS_PORT"),
		KeyID:           viper.GetString("KEY_ID"),
		AuthPageLink:    viper.GetString("AUTH_PAGE_LINK"),
		PrivateKeyBytes: []byte(viper.GetString("JWT_PRIVATE_KEY")),
		IDPConfFileName: viper.GetString("IDP_CONFIG_FILE_NAME"),
		TTL:             viper.GetDuration("ID_TOKEN_TTL"),
	}
}
