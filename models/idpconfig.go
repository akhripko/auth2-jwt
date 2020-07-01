package models

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type IDPType string
type ServiceType string

const (
	Microsoft IDPType = "microsoft"
	Google    IDPType = "google"
	Auth0     IDPType = "auth0"
	Github    IDPType = "github"

	StaffServiceType        ServiceType = "staff"
	ManningAgentServiceType ServiceType = "manning_agent"
)

type AuthConfig struct {
	Auth            *oauth2.Config
	Title           string
	IDPType         IDPType
	ServiceType     ServiceType
	Account         string
	EndPointProfile string
	Audience        string
	TTL             *time.Duration
}

type AccountConfig map[string]*AuthConfig

func (a AccountConfig) GetAuthConfig(key string) (*AuthConfig, error) {
	c, ok := a[key]
	if !ok {
		return nil, errors.Errorf("config was not found for `%s`", key)
	}
	return c, nil
}

func (a AccountConfig) GetOneAuthConfig() (*AuthConfig, error) {
	for _, authConf := range a {
		return authConf, nil
	}
	return nil, errors.New("account config is empty")
}

type IDPAccountConfigs map[string]AccountConfig

func (s IDPAccountConfigs) GetAccountConfig(key string) (AccountConfig, error) {
	c, ok := s[key]
	if !ok {
		return nil, errors.Errorf("account config was not found for '%s'", key)
	}
	if len(c) == 0 {
		return nil, errors.Errorf("account config is empty for '%s'", key)
	}
	return c, nil
}

func (s IDPAccountConfigs) SetAccountConfig(key string, config AccountConfig) {
	s[key] = config
}
