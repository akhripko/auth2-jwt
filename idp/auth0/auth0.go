package auth0

import (
	"encoding/json"
	"strings"

	"github.com/akhripko/auth2-jwt/models"
	"github.com/pkg/errors"
)

func BuildUser(data []byte, _ string, config *models.AuthConfig) (*models.User, error) {
	var values map[string]interface{}
	err := json.Unmarshal(data, &values)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build user model")
	}

	return &models.User{
		ID:         findID(values),
		Email:      getMapValue("email", values),
		Name:       getMapValue("name", values),
		NickName:   getMapValue("nickname", values),
		GivenName:  getMapValue("given_name", values),
		FamilyName: getMapValue("family_name", values),
		Picture:    getMapValue("picture", values),
		Context:    getMapValue("https://openocean.studio/scope", values),
	}, nil
}

func getMapValue(field string, values map[string]interface{}) string {
	value, ok := values[field]
	if !ok {
		return ""
	}

	if v, ok := value.(string); ok {
		return v
	}

	return ""
}

func findID(values map[string]interface{}) string {
	id := getMapValue("https://openocean.studio/federateduserid", values)
	if len(id) > 0 {
		return id
	}
	return extractID(getMapValue("sub", values))
}

func extractID(id string) string {
	ss := strings.Split(id, "|")
	if len(ss) == 0 {
		return id
	}
	return ss[len(ss)-1]
}
