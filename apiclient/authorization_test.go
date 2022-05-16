package apiclient

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApiKeyAuth(t *testing.T) {

	type args struct {
		key   string
		value string
	}

	tests := []struct {
		name string
		args args
		want *ApiClientAuthorization
	}{
		{
			"ApiKey authentication With Defined Key",
			args{
				"SomeApiKey",
				"SomeRandomKey",
			},
			&ApiClientAuthorization{
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
			&ApiClientAuthorization{
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

func TestNewStandardApiKey(t *testing.T) {
	// Arrange
	expected := ApiClientAuthorization{
		Key:   "ApiKey",
		Value: "someKey",
	}

	result := NewStandardApiKeyAuth("someKey")

	assert.Equalf(t, expected.String(), result.String(), "NewStandardApiKeyAuth() = %v, want %v", result.String(), expected.String())
}

func TestNewBearerTokenAuth(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want *ApiClientAuthorization
	}{
		{
			"Bearer token Authentication",
			args{
				token: "somebase64jwt",
			},
			&ApiClientAuthorization{
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

func TestNewBasicAuth(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name string
		args args
		want *ApiClientAuthorization
	}{
		{
			"Basic Authentication",
			args{
				username: "fakeUser",
				password: "fakePassword",
			},
			&ApiClientAuthorization{
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
