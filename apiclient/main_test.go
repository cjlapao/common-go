package apiclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultApiClient_SendRequestWithGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultApiClient{
		Client: server.Client(),
	}

	response, err := api.SendRequest(ApiClientOptions{
		Method:   GET,
		Protocol: "http",
		Host:     server.URL,
		Path:     "/foo/bar",
	})

	responseBodyRaw, errBody := ioutil.ReadAll(response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "ok", responseBody, "Response body should be ok")
}
