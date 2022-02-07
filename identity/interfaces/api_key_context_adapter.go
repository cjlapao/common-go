package interfaces

import (
	"errors"

	"github.com/cjlapao/common-go/identity/models"
)

var ErrApiKeyNotFound = errors.New("api key was not found")

type ApiKeyContextAdapter interface {
	Get(key string) (*models.ApiKey, error)
	Delete(key string) error
	Add(key *models.ApiKey) error
	Validate(key string, value string) (bool, error)
}
