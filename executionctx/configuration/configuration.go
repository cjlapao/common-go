package configuration

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/guard"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/security"
)

type ConfigurationProvider interface {
	UpsertKey(key string, value interface{}) error
	UpsertKeys(values map[string]interface{}) []error
	Get(key string) interface{}
	Clear(key string)
}

var configurationService *ConfigurationService

type ConfigurationService struct {
	Providers []ConfigurationProvider
}

func New() *ConfigurationService {
	if configurationService != nil {
		configurationService = nil
	}

	configurationService = &ConfigurationService{}
	if vault != nil {
		vault = nil
	}
	vault = make(map[string]interface{})

	return configurationService
}

func NewWithDefaults() *ConfigurationService {
	return New().RegisterDefaults()
}

func Get() *ConfigurationService {
	if configurationService != nil {
		return configurationService
	}

	return NewWithDefaults()
}

func (c *ConfigurationService) RegisterProvider(providers ...ConfigurationProvider) {
	for _, registerProvider := range providers {
		found := false
		for _, provider := range c.Providers {
			if reflect.TypeOf(provider) == reflect.TypeOf(registerProvider) {
				found = true
				break
			}
		}
		if !found {
			c.Providers = append(c.Providers, registerProvider)
		}
	}
}

func (c *ConfigurationService) UpsertKey(key string, value interface{}) error {
	if len(c.Providers) == 0 {
		return errors.New("no provider registered")
	}
	selectedProvider := c.Providers[0]
	for _, provider := range c.Providers {
		if reflect.TypeOf(provider) == reflect.TypeOf(CachedVaultConfigurationProvider{}) {
			selectedProvider = provider
			break
		}
	}

	return selectedProvider.UpsertKey(key, value)
}

func (c *ConfigurationService) UpsertKeys(values map[string]interface{}) []error {
	errorArray := make([]error, 0)
	if len(c.Providers) == 0 {
		error := errors.New("no provider registered")
		errorArray = append(errorArray, error)
	}
	return c.Providers[0].UpsertKeys(values)
}

func (c *ConfigurationService) Get(key string) interface{} {
	if len(c.Providers) == 0 {
		return nil
	}
	// Prioritizing the environment provider to get the keys
	for _, provider := range c.Providers {
		if reflect.TypeOf(provider) == reflect.TypeOf(EnvironmentConfigurationProvider{}) {
			result := provider.Get(key)
			if !guard.IsNill(result) {
				return result
			}
		}
	}

	for _, provider := range c.Providers {
		result := provider.Get(key)
		if !guard.IsNill(result) {
			return result
		}
	}

	return nil
}

func (c *ConfigurationService) GetString(key string) string {
	isKeyEmpty := guard.EmptyOrNil(key)
	if isKeyEmpty == nil {
		value := c.Get(key)
		if value == nil {
			return ""
		}
		return fmt.Sprint(value)
	}
	return ""
}

func (c *ConfigurationService) GetInt(key string) int {
	isKeyEmpty := guard.EmptyOrNil(key)
	if isKeyEmpty == nil {
		value, err := strconv.Atoi(fmt.Sprint(c.Get(key)))
		if err != nil {
			return 0
		}

		return value
	}

	return 0
}

func (c *ConfigurationService) GetBool(key string) bool {
	isKeyEmpty := guard.EmptyOrNil(key)
	if isKeyEmpty == nil {
		value, err := strconv.ParseBool(fmt.Sprint(c.Get(key)))
		if err != nil {
			return false
		}

		return value
	}

	return false
}

func (c *ConfigurationService) GetFloat(key string) float64 {
	isKeyEmpty := guard.EmptyOrNil(key)
	if isKeyEmpty == nil {
		value, err := strconv.ParseFloat(fmt.Sprint(c.Get(key)), 64)
		if err != nil {
			return value
		}

		return value
	}

	return 0
}

func (c *ConfigurationService) GetBase64(key string) string {
	isKeyEmpty := guard.EmptyOrNil(key)
	if isKeyEmpty != nil {
		return ""
	}

	value, err := security.DecodeBase64String(c.GetString(key))
	if err != nil {
		return ""
	}

	return value
}

func (c *ConfigurationService) Clear(key string) {
	for _, provider := range c.Providers {
		provider.Clear(key)
	}
}

func (c *ConfigurationService) LoadFromFile(path string) {
	if helper.FileExists(path) {
		content, err := helper.ReadFromFile(path)
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				index := strings.Index(line, "=")
				key := line[0:index]
				key = strings.TrimLeft(key, "\"")
				key = strings.TrimRight(key, "\"")

				value := line[index+1:]
				value = strings.TrimLeft(value, "\"")
				value = strings.TrimRight(value, "\"")
				c.UpsertKey(key, value)
			}
		}
	}
}
