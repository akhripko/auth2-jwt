package google

import (
	"encoding/json"

	"github.com/akhripko/auth2-jwt/models"
	"github.com/pkg/errors"
)

func BuildUser(data []byte, _ string, _ *models.AuthConfig) (*models.User, error) {
	var u User
	err := json.Unmarshal(data, &u)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build user model")
	}
	user := models.User(u)
	return &user, nil
}
