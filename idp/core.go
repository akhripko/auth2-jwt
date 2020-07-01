package idp

import (
	"github.com/akhripko/auth2-jwt/idp/auth0"
	"github.com/akhripko/auth2-jwt/idp/github"
	"github.com/akhripko/auth2-jwt/idp/google"
	"github.com/akhripko/auth2-jwt/idp/microsoft"
	"github.com/akhripko/auth2-jwt/models"
	"github.com/pkg/errors"
)

type UserBuilder func(data []byte, accessToken string, config *models.AuthConfig) (*models.User, error)

var userBuilders = map[models.IDPType]UserBuilder{
	models.Auth0:     auth0.BuildUser,
	models.Microsoft: microsoft.BuildUser,
	models.Google:    google.BuildUser,
	models.Github:    github.BuildUser,
}

func GetUserBuilder(key models.IDPType) (UserBuilder, error) {
	b, ok := userBuilders[key]
	if !ok {
		return nil, errors.Errorf("failed to get user builder: [%s] not supported IDP", key)
	}
	return b, nil
}
