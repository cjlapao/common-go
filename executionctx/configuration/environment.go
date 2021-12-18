package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/cjlapao/common-go/guard"
)

type EnvironmentConfigurationProvider struct {
}

func (ev EnvironmentConfigurationProvider) UpsertKey(key string, value interface{}) error {
	emptyKey := guard.EmptyOrNil(key, "key")
	emptyValue := guard.EmptyOrNil(value, "value")

	if emptyKey != nil {
		return emptyKey
	}

	if emptyValue != nil {
		return emptyValue
	}

	switch v := value.(type) {
	case int:
		os.Setenv(key, strconv.Itoa(v))
	case float64:
		os.Setenv(key, fmt.Sprintf("%v", v))
	case bool:
		os.Setenv(key, fmt.Sprintf("%v", v))
		os.Setenv(key, fmt.Sprintf("%v", v))
	case string:
		os.Setenv(key, v)
	default:
		jsonBytes, err := json.Marshal(v)
		if err == nil {
			os.Setenv(key, string(jsonBytes))
		}
	}

	return nil
}

func (ev EnvironmentConfigurationProvider) UpsertKeys(values map[string]interface{}) []error {
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

			ev.UpsertKey(key, value)
		}

		return errorArray
	}

	return nil
}

func (ev EnvironmentConfigurationProvider) Get(key string) interface{} {
	return os.Getenv(key)
}

func (ev EnvironmentConfigurationProvider) Clear(key string) {
	emptyKey := guard.EmptyOrNil(key, "key")

	if emptyKey == nil {
		os.Setenv(key, "")
	}
}
