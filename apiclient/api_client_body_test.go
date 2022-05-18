package apiclient

import (
	"strings"
	"testing"

	"github.com/cjlapao/common-go/helper"
	"github.com/stretchr/testify/assert"
)

func Test_ApiClientBody_NewApiClient_CreatesEmptyBody(t *testing.T) {

	body := NewApiClientBody()

	assert.NotNil(t, body)
	assert.Equalf(t, body.Type, BODY_TYPE_NONE, "body.Type = %v, want %v", body.Type, BODY_TYPE_NONE)
}

func Test_ApiClientBody_NewApiClientWithJsonBody_CreatesJsonBody(t *testing.T) {
	bodyContent := "{ \"someJson\": \"someData\" }"

	body := NewApiClientBody().Json([]byte(bodyContent))

	assert.NotNil(t, body)
	assert.Equalf(t, body.Type, BODY_TYPE_JSON, "body.Type = %v, want %v", body.Type, BODY_TYPE_JSON)
	assert.Equalf(t, bodyContent, string(body.Raw), "body.raw = %v, want %v", string(body.Raw), bodyContent)
}

func Test_ApiClientBody_NewApiClientWithTextBody_CreatesTextBody(t *testing.T) {
	bodyContent := "{ \"someJson\": \"someData\" }"

	body := NewApiClientBody().Text([]byte(bodyContent))

	assert.NotNil(t, body)
	assert.Equalf(t, body.Type, BODY_TYPE_TEXT, "body.Type = %v, want %v", body.Type, BODY_TYPE_TEXT)
	assert.Equalf(t, bodyContent, string(body.Raw), "body.raw = %v, want %v", string(body.Raw), bodyContent)
}

func Test_ApiClientBody_NewApiClientWithHtmlBody_CreatesHtmlBody(t *testing.T) {
	bodyContent := "{ \"someJson\": \"someData\" }"

	body := NewApiClientBody().Html([]byte(bodyContent))

	assert.NotNil(t, body)
	assert.Equalf(t, body.Type, BODY_TYPE_HTML, "body.Type = %v, want %v", body.Type, BODY_TYPE_HTML)
	assert.Equalf(t, bodyContent, string(body.Raw), "body.raw = %v, want %v", string(body.Raw), bodyContent)
}

func Test_ApiClientBody_NewApiClientWithUrlEncodedBody_CreatesUrlEncodedBody(t *testing.T) {
	body := NewApiClientBody().UrlEncoded().WithField("foo", "bar").WithField("bar", "foo")

	assert.NotNil(t, body)
	assert.Equalf(t, body.Type, BODY_TYPE_X_WWW_FORM_URLENCODED, "body.Type = %v, want %v", body.Type, BODY_TYPE_X_WWW_FORM_URLENCODED)
	assert.Equalf(t, 2, len(body.Fields), "body.Fields Count = %v, want %v", len(body.Fields), 2)
	assert.Equalf(t, "bar", body.Fields["foo"][0], "body.Fields.Foo = %v, want %v", body.Fields["foo"][0], "bar")
}

func Test_ApiClientBody_NewApiClientWithFormDataBody_CreatesFormDataBody(t *testing.T) {
	body := NewApiClientBody().FormData().WithField("foo", "bar").WithField("bar", "foo")

	assert.NotNil(t, body)
	assert.Equalf(t, body.Type, BODY_TYPE_FORM_DATA, "body.Type = %v, want %v", body.Type, BODY_TYPE_FORM_DATA)
	assert.Equalf(t, 2, len(body.Fields), "body.Fields Count = %v, want %v", len(body.Fields), 2)
	assert.Equalf(t, "bar", body.Fields["foo"][0], "body.Fields.Foo = %v, want %v", body.Fields["foo"][0], "bar")
}

func Test_ApiClientBody_WithFile_AddsCurrectFiles(t *testing.T) {
	filePath := "test.x"
	fileContent := "someContent"
	if helper.FileExists(filePath) {
		helper.DeleteFile(fileContent)
	}
	helper.WriteToFile(fileContent, filePath)

	body := NewApiClientBody().FormData().WithFile("file", filePath)

	helper.DeleteFile(filePath)
	assert.NotNil(t, body)
	assert.Equalf(t, 1, len(body.Files), "body.Files Count = %v, want %v", len(body.Files), 1)
	assert.Equal(t, len([]byte(fileContent)), len(body.Files["file"]["test.x"]))
}

func Test_ApiClientBody_WithField_AddsCurrectFields(t *testing.T) {
	body := NewApiClientBody().FormData().WithField("foo", "bar").WithField("bar", "foo")

	assert.NotNil(t, body)
	assert.Equalf(t, body.Type, BODY_TYPE_FORM_DATA, "body.Type = %v, want %v", body.Type, BODY_TYPE_FORM_DATA)
	assert.Equalf(t, 2, len(body.Fields), "body.Fields Count = %v, want %v", len(body.Fields), 2)
	assert.Equalf(t, "bar", body.Fields["foo"][0], "body.Fields.Foo = %v, want %v", body.Fields["foo"][0], "bar")
}

func Test_ApiClientBody_Get(t *testing.T) {
	formDataProcessedBody := NewApiClientBody().FormData().WithField("foo", "bar")
	formDataProcessedBody.Get()
	urlEncodedProcessedBody := NewApiClientBody().UrlEncoded().WithField("foo", "bar")
	urlEncodedProcessedBody.Get()

	tests := []struct {
		name           string
		body           *ApiClientBody
		wantBody       string
		containsInBody string
	}{
		{
			name:     "Json",
			body:     NewApiClientBody().Json([]byte("{}")),
			wantBody: "{}",
		},
		{
			name:     "Text",
			body:     NewApiClientBody().Text([]byte("{text}")),
			wantBody: "{text}",
		},
		{
			name:     "Html",
			body:     NewApiClientBody().Html([]byte("{html}")),
			wantBody: "{html}",
		},
		{
			name:           "FormData",
			body:           formDataProcessedBody,
			containsInBody: "Content-Disposition: form-data; name=\"foo\"\r\n\r\nbar\r\n",
		},
		{
			name:     "Url Encoded",
			body:     urlEncodedProcessedBody,
			wantBody: "foo=bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := tt.body.Get()

			if tt.containsInBody != "" {
				if !strings.ContainsAny(body.String(), tt.containsInBody) {
					t.Errorf("body.Value = %v, want %v", body.String(), tt.containsInBody)
				}
			}
			if tt.wantBody != "" {
				assert.Equalf(t, tt.wantBody, body.String(), "body.Files Count = %v, want %v", body.String(), tt.wantBody)
			}
		})
	}
}

func Test_ApiClientBody_GetHeader(t *testing.T) {
	formDataProcessedBody := NewApiClientBody().FormData().WithField("foo", "bar")
	formDataProcessedBody.Get()
	urlEncodedProcessedBody := NewApiClientBody().UrlEncoded().WithField("foo", "bar")
	urlEncodedProcessedBody.Get()

	tests := []struct {
		name      string
		body      *ApiClientBody
		wantKey   string
		wantValue string
	}{
		{
			name:      "Json",
			body:      NewApiClientBody().Json([]byte("{}")),
			wantKey:   "Content-Type",
			wantValue: "application/json;charset=UTF-8",
		},
		{
			name:      "Text",
			body:      NewApiClientBody().Text([]byte("{}")),
			wantKey:   "Content-Type",
			wantValue: "plain/text",
		},
		{
			name:      "Html",
			body:      NewApiClientBody().Html([]byte("{}")),
			wantKey:   "Content-Type",
			wantValue: "text/html",
		},
		{
			name:      "FormData",
			body:      formDataProcessedBody,
			wantKey:   "Content-Type",
			wantValue: "multipart/form-data; boundary=",
		},
		{
			name:      "Url Encoded",
			body:      urlEncodedProcessedBody,
			wantKey:   "Content-Type",
			wantValue: "application/x-www-form-urlencoded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, value := tt.body.GetHeader()

			assert.Equalf(t, tt.wantKey, key, "body.Header.Key = %v, want %v", tt.wantKey, key)
			if !strings.HasPrefix(value, tt.wantValue) {
				t.Errorf("body.Header.Value = %v, want %v", tt.wantValue, value)
			}
		})
	}
}
