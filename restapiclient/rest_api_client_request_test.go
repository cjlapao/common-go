package restapiclient

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RestApiClientRequest_ParseURL(t *testing.T) {
	tests := []struct {
		name                    string
		url                     string
		expectedHost            string
		expectedScheme          string
		expectedPath            string
		expectedQueryParameters url.Values
		wantError               bool
	}{
		{
			name:                    "good host",
			url:                     "http://example.com",
			expectedHost:            "example.com",
			expectedScheme:          "http",
			expectedPath:            "",
			expectedQueryParameters: url.Values{},
			wantError:               false,
		},
		{
			name:                    "good host",
			url:                     "http://example.com/foo",
			expectedHost:            "example.com",
			expectedScheme:          "http",
			expectedPath:            "/foo",
			expectedQueryParameters: url.Values{},
			wantError:               false,
		},
		{
			name:           "good host",
			url:            "http://example.com/foo?test=true",
			expectedHost:   "example.com",
			expectedScheme: "http",
			expectedPath:   "/foo",
			expectedQueryParameters: url.Values{
				"test": []string{"true"},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := RestApiClientRequest{}
			err := request.ParseUrl(tt.url)
			if err != nil && !tt.wantError {
				t.Error("expected error to be nil but found error")
			}
			if tt.wantError && err == nil {
				t.Error("expected error not to be nil but found no error")
			}

			assert.Equal(t, tt.expectedHost, request.URL.Host)
			assert.Equal(t, tt.expectedScheme, request.URL.Scheme)
			assert.Equal(t, tt.expectedPath, request.URL.Path)
			for wantKey, wantValue := range tt.expectedQueryParameters {
				found := false
				for gotKey, gotValue := range request.URL.Query() {
					if gotKey == wantKey && gotValue[0] == wantValue[0] {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%v = %v, not found", wantKey, wantValue)
				}
			}
		})
	}
}

func Test_RestApiClientRequest_AddHeader_WithNoCtor(t *testing.T) {
	request := RestApiClientRequest{}

	request.AddHeader("X-Requested-With", "some text")

	assert.NotNil(t, request.Headers)
	assert.Equal(t, 1, len(request.Headers))
	assert.Equal(t, "some text", request.Headers["X-Requested-With"])
}

func Test_RestApiClientRequest_AddHeader_WithExistingHeader_GetsAppended(t *testing.T) {
	request := RestApiClientRequest{}
	request.Headers = map[string]string{
		"Content-Length": "10",
	}

	request.AddHeader("X-Requested-With", "some text")

	assert.NotNil(t, request.Headers)
	assert.Equal(t, 2, len(request.Headers))
	assert.Equal(t, "some text", request.Headers["X-Requested-With"])
}
