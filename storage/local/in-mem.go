package local

import (
	"github.com/akhripko/auth2-jwt/models"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

func NewInMemStorage(data models.IDPAccountConfigs) *InMemStorage {
	return &InMemStorage{
		data: data,
	}
}

func BuildInMemStorageFromYaml(yamlBytes []byte) (*InMemStorage, error) {
	data := make(models.IDPAccountConfigs)
	err := yaml.Unmarshal(yamlBytes, data)
	if err != nil {
		return nil, errors.Wrap(err, "build from yaml failed ")
	}
	return &InMemStorage{
		data: data,
	}, nil
}

type InMemStorage struct {
	data models.IDPAccountConfigs
}

func (s *InMemStorage) Check() error {
	if s.data == nil {
		return errors.New("local storage was not initialized")
	}
	return nil
}

func (s *InMemStorage) ReadAccountConfig(account string) (models.AccountConfig, error) {
	if err := s.Check(); err != nil {
		return nil, errors.Wrap(err, "failed to read account data")
	}
	return s.data.GetAccountConfig(account)
}
