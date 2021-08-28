package executionctx

import (
	"errors"
	"fmt"
	"os"

	"github.com/cjlapao/common-go/guard"
)

var vault = make(map[string]interface{})
var configurationService *Configuration

type Configuration struct {
}

func NewConfigService() *Configuration {
	if configurationService != nil {
		configurationService = nil
	}

	configurationService = &Configuration{}
	if vault != nil {
		vault = nil
	}
	vault = make(map[string]interface{})

	return configurationService
}

func GetConfigService() *Configuration {
	if configurationService != nil {
		return configurationService
	}

	return NewConfigService()
}

func (c *Configuration) UpsertKey(key string, value interface{}) error {
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

func (c *Configuration) UpsertKeys(values map[string]interface{}) []error {
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

func (c *Configuration) Get(key string) interface{} {
	var value interface{}
	value = c.getFromEnvironment(key)
	if value == "" {
		value = c.getFromVault(key)
	}

	return value
}

func (c *Configuration) Clear() {
	vault = make(map[string]interface{})
}

func (c *Configuration) getFromVault(key string) interface{} {
	for keyIndex, keyValue := range vault {
		if key == keyIndex {
			fmt.Print("found")
			return keyValue
		}
	}

	return nil
}

func (c *Configuration) getFromEnvironment(key string) string {
	return os.Getenv(key)
}
