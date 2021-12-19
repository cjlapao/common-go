package configuration

import (
	"errors"

	"github.com/cjlapao/common-go/guard"
)

var vault = make(map[string]interface{})

type CachedVaultConfigurationProvider struct {
}

func (ev CachedVaultConfigurationProvider) UpsertKey(key string, value interface{}) error {
	emptyKey := guard.EmptyOrNil(key, "key")
	emptyValue := guard.EmptyOrNil(value, "value")

	if emptyKey != nil {
		return emptyKey
	}

	if emptyValue != nil {
		return emptyValue
	}

	vault[key] = value

	return nil
}

func (ev CachedVaultConfigurationProvider) UpsertKeys(values map[string]interface{}) []error {
	errorArray := make([]error, 0)

	if values == nil {
		errorArray = append(errorArray, errors.New("array is nil"))
		return errorArray
	}

	if len(values) > 0 {
		for key, value := range values {
			emptykey := guard.EmptyOrNil(key)
			emptyValue := guard.EmptyOrNil(value)

			if emptykey != nil {
				errorArray = append(errorArray, emptykey)
			}

			if emptyValue != nil {
				errorArray = append(errorArray, emptyValue)
			}

			vault[key] = value
		}
		return errorArray
	}

	return nil
}

func (ev CachedVaultConfigurationProvider) Get(key string) interface{} {
	for keyIndex, keyValue := range vault {
		if key == keyIndex {
			return keyValue
		}
	}

	return nil
}

func (ev CachedVaultConfigurationProvider) Clear(key string) {
	emptyKey := guard.EmptyOrNil(key, "key")

	if emptyKey == nil {
		delete(vault, key)
		//TODO: Implement the logic to remove it from the cached vault
	}
}
