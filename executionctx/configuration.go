package executionctx

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

var vault = make(map[string]interface{})
var configurationService *Configuration

type Configuration struct {
}

func GetConfigurationProvider() *Configuration {
	if configurationService != nil {
		return configurationService
	}
	configurationService = &Configuration{}
	vault["port"] = 500
	vault["url"] = "hostname"
	return configurationService
}

func (c *Configuration) UpsertKey(key string, value interface{}) error {
	if value == nil {
		return errors.New("Value cannot be empty")
	}

	reflect.Zero(reflect.TypeOf(value)).Interface()
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
