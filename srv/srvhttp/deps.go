package srvhttp

import "github.com/akhripko/auth2-jwt/models"

type IDPConfigsStorage interface {
	Check() error
	ReadAccountConfig(account string) (models.AccountConfig, error)
}
