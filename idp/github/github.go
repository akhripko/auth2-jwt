package github

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/akhripko/auth2-jwt/models"
	"github.com/pkg/errors"
)

func BuildUser(data []byte, accessToken string, _ *models.AuthConfig) (*models.User, error) {
	var u User
	err := json.Unmarshal(data, &u)
	if err != nil {
		return nil, errors.Wrap(err, "[github] failed to build user model")
	}

	teams, err := fetchTeams(accessToken)
	if err != nil {
		return nil, errors.Wrap(err, "[github] failed to read user teams")
	}

	user := models.User{
		ID:       strconv.Itoa(int(u.ID)),
		Name:     u.Name,
		NickName: u.NickName,
		Picture:  u.Picture,
		Groups:   teams,

		// note: email can be empty if user hasn't public email
		Email: u.Email,

		// note: github doesn't provide these fields
		GivenName:  "",
		FamilyName: "",
	}

	return &user, nil
}

func fetchTeams(accessToken string) (map[string]string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/teams", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var teams []Team
	err = json.Unmarshal(data, &teams)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(teams))
	for _, team := range teams {
		result[team.Organization.Name+"/"+team.Name] = strconv.Itoa(int(team.ID))
	}

	return result, nil
}
