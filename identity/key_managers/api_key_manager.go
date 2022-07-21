package key_managers

import (
	"strings"

	cryptorand "github.com/cjlapao/common-go-cryptorand"
	"github.com/cjlapao/common-go/constants"
	"github.com/cjlapao/common-go/guard"
	"github.com/cjlapao/common-go/identity/interfaces"
	"github.com/cjlapao/common-go/identity/models"
)

var globalApiKeyManager *ApiKeyManager

type ApiKeyManager struct {
	ContextAdapter interfaces.ApiKeyContextAdapter
	Keys           []*models.ApiKey
}

func EmptyApiKeyManager() *ApiKeyManager {
	apiKeyManager := ApiKeyManager{
		Keys: make([]*models.ApiKey, 0),
	}

	globalApiKeyManager = &apiKeyManager

	return globalApiKeyManager
}

func GetApiKeyManager() *ApiKeyManager {
	if globalApiKeyManager == nil {
		EmptyApiKeyManager()
	}

	return globalApiKeyManager
}

func (apiKeyManager *ApiKeyManager) Get(key string) (*models.ApiKey, error) {
	for _, apiKey := range apiKeyManager.Keys {
		if strings.EqualFold(apiKey.Name, key) {
			return apiKey, nil
		}
	}

	dbKey, err := apiKeyManager.ContextAdapter.Get(key)
	if err != nil {
		return nil, err
	}

	return dbKey, nil
}

func (apiKeyManager *ApiKeyManager) Add(key string, value string) error {
	if err := guard.EmptyOrNil(key); err != nil {
		return err
	}

	if err := guard.EmptyOrNil(value); err != nil {
		return err
	}

	var cachedApiKey *models.ApiKey

	for idx, apiKey := range apiKeyManager.Keys {
		if strings.EqualFold(apiKey.Name, key) {
			apiKeyManager.Keys[idx].KeyValue = value
			cachedApiKey = apiKey
			break
		}
	}

	if cachedApiKey == nil {
		cachedApiKey = &models.ApiKey{
			ID:       cryptorand.GetRandomString(constants.ID_SIZE),
			Name:     key,
			KeyValue: value,
		}
	}

	if err := apiKeyManager.ContextAdapter.Add(cachedApiKey); err != nil {
		return err
	}

	return nil
}
