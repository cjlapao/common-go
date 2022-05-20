package restapiclient

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RestApiClientAuthorization_NewApiKeyAuth(t *testing.T) {

	type args struct {
		key   string
		value string
	}

	tests := []struct {
		name string
		args args
		want *RestApiClientAuthorization
	}{
		{
			"ApiKey authentication With Defined Key",
			args{
				"SomeApiKey",
				"SomeRandomKey",
			},
			&RestApiClientAuthorization{
				"SomeApiKey",
				"SomeRandomKey",
			},
		},
		{
			"ApiKey authentication With no Key defined",
			args{
				"",
				"SomeRandomKey",
			},
			&RestApiClientAuthorization{
				"ApiKey",
				"SomeRandomKey",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewApiKeyAuth(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApiKeyAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RestApiClientAuthorization_NewStandardApiKey(t *testing.T) {
	// Arrange
	expected := RestApiClientAuthorization{
		Key:   "ApiKey",
		Value: "someKey",
	}

	result := NewStandardApiKeyAuth("someKey")

	assert.Equalf(t, expected.String(), result.String(), "NewStandardApiKeyAuth() = %v, want %v", result.String(), expected.String())
}

func Test_RestApiClientAuthorization_NewBearerTokenAuth(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want *RestApiClientAuthorization
	}{
		{
			"Bearer token Authentication",
			args{
				token: "somebase64jwt",
			},
			&RestApiClientAuthorization{
				Key:   "Bearer",
				Value: "somebase64jwt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBearerTokenAuth(tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBearerTokenAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RestApiClientAuthorization_NewBasicAuth(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name string
		args args
		want *RestApiClientAuthorization
	}{
		{
			"Basic Authentication",
			args{
				username: "fakeUser",
				password: "fakePassword",
			},
			&RestApiClientAuthorization{
				Key:   "Basic",
				Value: "ZmFrZVVzZXI6ZmFrZVBhc3N3b3Jk",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBasicAuth(tt.args.username, tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBasicAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test__RestApiClientAuthorization_GetHeader(t *testing.T) {
	auth := NewBearerTokenAuth("testToken")

	key, value := auth.GetHeader()

	assert.Equalf(t, "Authorization", key, "Key Header = %v, want \"Authorization\"", key)
	assert.Equal(t, 1, len(value))
	assert.Equal(t, "Bearer testToken", value[0])
}
