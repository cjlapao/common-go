package apiclient

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func Test_ApiBodyTest_GetHeader(t *testing.T) {
	tests := []struct {
		name      string
		bodyType  ApiClientBodyType
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

			assert.Equalf(t, tt.wantKey, key, "ApiClientBodyType.Key = %v, want %v", tt.wantKey, key)
			assert.Equalf(t, tt.wantValue, value, "ApiClientBodyType.Value = %v, want %v", tt.wantValue, value)
		})
	}
}

func Test_ApiBodyTest_UnmarshalJson(t *testing.T) {
	testObj := struct {
		SomeType ApiClientBodyType `json:"someType"`
	}{}

	plainJson := "{ \"someType\": \"FORM-DATA\" }"

	json.Unmarshal([]byte(plainJson), &testObj)

	assert.Equalf(t, testObj.SomeType, BODY_TYPE_FORM_DATA, "testObj.SomeType = %v, want %v", testObj.SomeType, BODY_TYPE_FORM_DATA)
}

func Test_ApiBodyTest_MarshalJson(t *testing.T) {
	testObj := struct {
		SomeType ApiClientBodyType `json:"someType"`
	}{
		SomeType: BODY_TYPE_FORM_DATA,
	}

	expectedPlainJson := "{\"someType\":\"FORM-DATA\"}"
	plainJson, err := json.Marshal(testObj)

	assert.Nil(t, err)
	assert.Equalf(t, string(plainJson), expectedPlainJson, "json = %v, want %v", string(plainJson), expectedPlainJson)
}

func Test_ApiBodyTest_UnmarshalYaml(t *testing.T) {
	testObj := struct {
		SomeType ApiClientBodyType `yaml:"someType"`
	}{}

	plainYaml := "someType: FORM-DATA"

	err := yaml.Unmarshal([]byte(plainYaml), &testObj)
	assert.Nil(t, err)

	assert.Equalf(t, testObj.SomeType, BODY_TYPE_FORM_DATA, "testObj.SomeType = %v, want %v", testObj.SomeType, BODY_TYPE_FORM_DATA)
}

func Test_ApiBodyTest_MarshalYaml(t *testing.T) {
	testObj := struct {
		SomeType ApiClientBodyType `yaml:"someType"`
	}{
		SomeType: BODY_TYPE_FORM_DATA,
	}

	expectedPlainYaml := "someType: FORM-DATA\n"
	plainYaml, err := yaml.Marshal(testObj)

	assert.Nil(t, err)
	assert.Equalf(t, string(plainYaml), expectedPlainYaml, "yaml = %v, want %v", string(plainYaml), expectedPlainYaml)
}
