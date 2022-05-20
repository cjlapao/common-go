package restapiclient

import (
	"fmt"

	"github.com/cjlapao/common-go/security"
)

type RestApiClientAuthorization struct {
	Key   string
	Value string
}

func NewApiKeyAuth(key string, value string) *RestApiClientAuthorization {
	if key == "" {
		key = "ApiKey"
	}

	return &RestApiClientAuthorization{
		Key:   key,
		Value: value,
	}
}

func NewStandardApiKeyAuth(value string) *RestApiClientAuthorization {
	return NewApiKeyAuth("ApiKey", value)
}

func NewBearerTokenAuth(token string) *RestApiClientAuthorization {
	return &RestApiClientAuthorization{
		Key:   "Bearer",
		Value: token,
	}
}

func NewBasicAuth(username string, password string) *RestApiClientAuthorization {
	value, err := security.EncodeString(fmt.Sprintf("%v:%v", username, password))

	if err != nil {
		value = ""
	}

	return &RestApiClientAuthorization{
		Key:   "Basic",
		Value: value,
	}
}

func (auth *RestApiClientAuthorization) String() string {
	return fmt.Sprintf("%v %v", auth.Key, auth.Value)
}

func (auth *RestApiClientAuthorization) GetHeader() (key string, value []string) {
	return "Authorization", []string{auth.String()}
}
