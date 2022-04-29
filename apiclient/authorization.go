package apiclient

import (
	"fmt"

	"github.com/cjlapao/common-go/security"
)

type ApiClientAuthorization struct {
	Key   string
	Value string
}

func NewApiKeyAuth(key string, value string) *ApiClientAuthorization {
	if key == "" {
		key = "ApiKey"
	}

	return &ApiClientAuthorization{
		Key:   key,
		Value: value,
	}
}

func NewBearerTokenAuth(token string) *ApiClientAuthorization {
	return &ApiClientAuthorization{
		Key:   "Bearer",
		Value: token,
	}
}

func NewBasicAuth(username string, password string) *ApiClientAuthorization {
	value, err := security.EncodeString(fmt.Sprintf("%v:%v", username, password))

	if err != nil {
		value = ""
	}

	return &ApiClientAuthorization{
		Key:   "Basic",
		Value: value,
	}
}

func (auth *ApiClientAuthorization) String() string {
	return fmt.Sprintf("%v %v", auth.Key, auth.Value)
}
