package restapiclient

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func Test_RestApiClientMethod_GetHeader(t *testing.T) {
	tests := []struct {
		name      string
		apiMethod RestApiClientMethod
		wantValue string
	}{
		{
			apiMethod: API_METHOD_GET,
			wantValue: "GET",
		},
		{
			apiMethod: API_METHOD_POST,
			wantValue: "POST",
		},
		{
			apiMethod: API_METHOD_PUT,
			wantValue: "PUT",
		},
		{
			apiMethod: API_METHOD_PATCH,
			wantValue: "PATCH",
		},
		{
			apiMethod: API_METHOD_DELETE,
			wantValue: "DELETE",
		},
		{
			apiMethod: API_METHOD_CONNECT,
			wantValue: "CONNECT",
		},
		{
			apiMethod: API_METHOD_HEAD,
			wantValue: "HEAD",
		},
		{
			apiMethod: API_METHOD_OPTIONS,
			wantValue: "OPTIONS",
		},
		{
			apiMethod: API_METHOD_TRACE,
			wantValue: "TRACE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := tt.apiMethod.String()

			assert.Equalf(t, tt.wantValue, value, "RestApiClientMethod() = %v, want %v", value, tt.wantValue)
		})
	}
}

func Test_RestApiClientMethod_UnmarshalJson(t *testing.T) {
	testObj := struct {
		SomeMethod RestApiClientMethod `json:"someMethod"`
	}{}

	plainJson := "{ \"someMethod\": \"GET\" }"

	json.Unmarshal([]byte(plainJson), &testObj)

	assert.Equalf(t, API_METHOD_GET, testObj.SomeMethod, "testObj.SomeMethod = %v, want %v", testObj.SomeMethod, API_METHOD_GET)
}

func Test_RestApiClientMethod_MarshalJson(t *testing.T) {
	testObj := struct {
		SomeMethod RestApiClientMethod `json:"someMethod"`
	}{
		SomeMethod: API_METHOD_OPTIONS,
	}

	expectedPlainJson := "{\"someMethod\":\"OPTIONS\"}"
	plainJson, err := json.Marshal(testObj)

	assert.Nil(t, err)
	assert.Equalf(t, expectedPlainJson, string(plainJson), "json = %v, want %v", string(plainJson), expectedPlainJson)
}

func Test_RestApiClientMethod_UnmarshalYaml(t *testing.T) {
	testObj := struct {
		SomeMethod RestApiClientMethod `yaml:"someMethod"`
	}{}

	plainYaml := "someMethod: PUT"

	err := yaml.Unmarshal([]byte(plainYaml), &testObj)
	assert.Nil(t, err)

	assert.Equalf(t, API_METHOD_PUT, testObj.SomeMethod, "testObj.SomeMethod = %v, want %v", testObj.SomeMethod, API_METHOD_PUT)
}

func Test_RestApiClientMethod_MarshalYaml(t *testing.T) {
	testObj := struct {
		SomeMethod RestApiClientMethod `yaml:"someMethod"`
	}{
		SomeMethod: API_METHOD_CONNECT,
	}

	expectedPlainYaml := "someMethod: CONNECT\n"
	plainYaml, err := yaml.Marshal(testObj)

	assert.Nil(t, err)
	assert.Equalf(t, expectedPlainYaml, string(plainYaml), "yaml = %v, want %v", string(plainYaml), expectedPlainYaml)
}

func Test_RestApiClientMethod_FromString(t *testing.T) {
	var method RestApiClientMethod
	methodString := method.FromString("TRACE")

	assert.Equalf(t, API_METHOD_TRACE, methodString, "RestApiClientMethod() = %v, want %v", methodString, API_METHOD_GET)
}
