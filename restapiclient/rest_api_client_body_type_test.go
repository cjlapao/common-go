package restapiclient

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func Test_RestApiClientBodyType_GetHeader(t *testing.T) {
	tests := []struct {
		name      string
		bodyType  RestApiClientBodyType
		wantKey   string
		wantValue string
	}{
		{
			bodyType:  BODY_TYPE_JSON,
			wantKey:   "Content-Type",
			wantValue: "application/json;charset=UTF-8",
		},
		{
			bodyType:  BODY_TYPE_TEXT,
			wantKey:   "Content-Type",
			wantValue: "plain/text",
		},
		{
			bodyType:  BODY_TYPE_HTML,
			wantKey:   "Content-Type",
			wantValue: "text/html",
		},
		{
			bodyType:  BODY_TYPE_FORM_DATA,
			wantKey:   "Content-Type",
			wantValue: "multipart/form-data",
		},
		{
			bodyType:  BODY_TYPE_X_WWW_FORM_URLENCODED,
			wantKey:   "Content-Type",
			wantValue: "application/x-www-form-urlencoded",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, value := tt.bodyType.GetHeader()

			assert.Equalf(t, tt.wantKey, key, "RestApiClientBodyType.Key = %v, want %v", tt.wantKey, key)
			assert.Equalf(t, tt.wantValue, value, "RestApiClientBodyType.Value = %v, want %v", tt.wantValue, value)
		})
	}
}

func Test_RestApiClientBodyType_UnmarshalJson(t *testing.T) {
	testObj := struct {
		SomeType RestApiClientBodyType `json:"someType"`
	}{}

	plainJson := "{ \"someType\": \"FORM-DATA\" }"

	json.Unmarshal([]byte(plainJson), &testObj)

	assert.Equalf(t, testObj.SomeType, BODY_TYPE_FORM_DATA, "testObj.SomeType = %v, want %v", testObj.SomeType, BODY_TYPE_FORM_DATA)
}

func Test_RestApiClientBodyType_MarshalJson(t *testing.T) {
	testObj := struct {
		SomeType RestApiClientBodyType `json:"someType"`
	}{
		SomeType: BODY_TYPE_FORM_DATA,
	}

	expectedPlainJson := "{\"someType\":\"FORM-DATA\"}"
	plainJson, err := json.Marshal(testObj)

	assert.Nil(t, err)
	assert.Equalf(t, string(plainJson), expectedPlainJson, "json = %v, want %v", string(plainJson), expectedPlainJson)
}

func Test_RestApiClientBodyType_UnmarshalYaml(t *testing.T) {
	testObj := struct {
		SomeType RestApiClientBodyType `yaml:"someType"`
	}{}

	plainYaml := "someType: FORM-DATA"

	err := yaml.Unmarshal([]byte(plainYaml), &testObj)
	assert.Nil(t, err)

	assert.Equalf(t, testObj.SomeType, BODY_TYPE_FORM_DATA, "testObj.SomeType = %v, want %v", testObj.SomeType, BODY_TYPE_FORM_DATA)
}

func Test_RestApiClientBodyType_MarshalYaml(t *testing.T) {
	testObj := struct {
		SomeType RestApiClientBodyType `yaml:"someType"`
	}{
		SomeType: BODY_TYPE_FORM_DATA,
	}

	expectedPlainYaml := "someType: FORM-DATA\n"
	plainYaml, err := yaml.Marshal(testObj)

	assert.Nil(t, err)
	assert.Equalf(t, string(plainYaml), expectedPlainYaml, "yaml = %v, want %v", string(plainYaml), expectedPlainYaml)
}

func Test_RestApiClientBodyType_FromString(t *testing.T) {
	var bodyType RestApiClientBodyType
	bodyTypeString := bodyType.FromString("JSON")

	assert.Equalf(t, BODY_TYPE_JSON, bodyTypeString, "RestApiClientBodyType() = %v, want %v", bodyTypeString, BODY_TYPE_JSON)
}
