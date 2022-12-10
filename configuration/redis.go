package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/cjlapao/common-go/guard"
)

type RedisConfigurationProvider struct {
	connectionString string
	client           string
}

func NewRedisConfigurationProvider(connString string) *RedisConfigurationProvider {
	result := RedisConfigurationProvider{
		connectionString: connString,
	}

	return &result
}

func (ev RedisConfigurationProvider) UpsertKey(key string, value interface{}) error {
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

func (ev RedisConfigurationProvider) UpsertKeys(values map[string]interface{}) []error {
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

func (ev RedisConfigurationProvider) Get(key string) interface{} {
	return os.Getenv(key)
}

func (ev RedisConfigurationProvider) Clear(key string) {
	emptyKey := guard.EmptyOrNil(key, "key")

	if emptyKey == nil {
		os.Setenv(key, "")
	}
}
