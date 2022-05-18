package apiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/security"
	"github.com/stretchr/testify/assert"
)

type TestBody struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

func TestDefaultApiClient_SendRequestWithGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "GET", r.Method, "Expected to be GET method")
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
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func TestDefaultApiClient_SendRequestWithGetAndTokenAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "GET", r.Method, "Expected to be GET method")
		token := r.Header.Get("Authorization")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		assert.NotNilf(t, token, "token should not be nil was %v", token)
		assert.Equalf(t, "Bearer abc", token, "Token should be \"Bearer abc\" found %v", token)
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultApiClient{
		Client: server.Client(),
	}

	response, err := api.SendRequest(ApiClientOptions{
		Method:        GET,
		Protocol:      "http",
		Host:          server.URL,
		Path:          "/foo/bar",
		Authorization: *NewBearerTokenAuth("abc"),
	})

	responseBodyRaw, errBody := ioutil.ReadAll(response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func TestDefaultApiClient_SendRequestWithGetAndBasicAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "GET", r.Method, "Expected to be GET method")
		token := r.Header.Get("Authorization")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		assert.NotNilf(t, token, "token should not be nil was %v", token)
		user, _ := security.EncodeString("testUser:testPassword")
		assert.Equalf(t, fmt.Sprintf("Basic %v", user), token, "Token should be \"%v\" found %v", user, token)
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultApiClient{
		Client: server.Client(),
	}

	response, err := api.SendRequest(ApiClientOptions{
		Method:        GET,
		Protocol:      "http",
		Host:          server.URL,
		Path:          "/foo/bar",
		Authorization: *NewBasicAuth("testUser", "testPassword"),
	})

	responseBodyRaw, errBody := ioutil.ReadAll(response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func TestDefaultApiClient_SendRequestWithGetAndApiKeyAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "GET", r.Method, "Expected to be GET method")
		apiKey := r.Header.Get("Authorization")
		expectedApiKey := "TestKey someKey"
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		assert.NotNilf(t, apiKey, "token should not be nil was %v", apiKey)
		assert.Equalf(t, expectedApiKey, apiKey, "Token should be \"%v\" found %v", expectedApiKey, apiKey)
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultApiClient{
		Client: server.Client(),
	}

	response, err := api.SendRequest(ApiClientOptions{
		Method:        GET,
		Protocol:      "http",
		Host:          server.URL,
		Path:          "/foo/bar",
		Authorization: *NewApiKeyAuth("TestKey", "someKey"),
	})

	responseBodyRaw, errBody := ioutil.ReadAll(response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func TestDefaultApiClient_SendRequestWithGetAndStandardApiKeyAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "GET", r.Method, "Expected to be GET method")
		apiKey := r.Header.Get("Authorization")
		expectedApiKey := "ApiKey someKey"
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		assert.NotNilf(t, apiKey, "token should not be nil was %v", apiKey)
		assert.Equalf(t, expectedApiKey, apiKey, "Token should be \"%v\" found %v", expectedApiKey, apiKey)
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultApiClient{
		Client: server.Client(),
	}

	response, err := api.SendRequest(ApiClientOptions{
		Method:        GET,
		Protocol:      "http",
		Host:          server.URL,
		Path:          "/foo/bar",
		Authorization: *NewStandardApiKeyAuth("someKey"),
	})

	responseBodyRaw, errBody := ioutil.ReadAll(response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func TestDefaultApiClient_SendRequestWithPostAndJsonBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotNilf(t, r.Body, "body should not be nil")
		var body TestBody
		json.NewDecoder(r.Body).Decode(&body)

		assert.Equalf(t, "POST", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")

		assert.Equalf(t, "testUser", body.Username, "username = %v, want %v", body.Username, "testUser")
		assert.Equalf(t, "testPassword", body.Password, "password = %v, want %v", body.Password, "testPassword")
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultApiClient{
		Client: server.Client(),
	}

	testObj := TestBody{
		Username: "testUser",
		Password: "testPassword",
	}

	marshaledtestObj, _ := json.Marshal(testObj)

	response, err := api.SendRequest(ApiClientOptions{
		Method:   POST,
		Protocol: "http",
		Host:     server.URL,
		Path:     "/foo/bar",
		Body:     *NewApiClientBody().Json(marshaledtestObj),
	})

	responseBodyRaw, errBody := ioutil.ReadAll(response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func TestDefaultApiClient_SendRequestWithPostAndXUrlEncodedBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotNilf(t, r.Body, "body should not be nil")
		var body TestBody
		pError := r.ParseForm()
		assert.Nilf(t, pError, "parsing form should not contain errors")

		body.Password = r.Form.Get("password")
		body.Username = r.Form.Get("username")

		assert.Equalf(t, "POST", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")

		assert.Equalf(t, "testUser", body.Username, "username = %v, want %v", body.Username, "testUser")
		assert.Equalf(t, "testPassword", body.Password, "password = %v, want %v", body.Password, "testPassword")
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultApiClient{
		Client: server.Client(),
	}

	response, err := api.SendRequest(ApiClientOptions{
		Method:   POST,
		Protocol: "http",
		Host:     server.URL,
		Path:     "/foo/bar",
		Body:     *NewApiClientBody().UrlEncoded().WithField("username", "testUser").WithField("password", "testPassword"),
	})

	responseBodyRaw, errBody := ioutil.ReadAll(response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func TestDefaultApiClient_SendRequestWithPostAndFormDataBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotNilf(t, r.Body, "body should not be nil")
		var body TestBody
		pError := r.ParseMultipartForm(0)
		assert.Nilf(t, pError, "parsing form should not contain errors")

		body.Password = r.Form.Get("password")
		body.Username = r.Form.Get("username")

		assert.Equalf(t, "POST", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")

		assert.Equalf(t, "testUser", body.Username, "username = %v, want %v", body.Username, "testUser")
		assert.Equalf(t, "testPassword", body.Password, "password = %v, want %v", body.Password, "testPassword")
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultApiClient{
		Client: server.Client(),
	}

	response, err := api.SendRequest(ApiClientOptions{
		Method:   POST,
		Protocol: "http",
		Host:     server.URL,
		Path:     "/foo/bar",
		Body:     *NewApiClientBody().FormData().WithField("username", "testUser").WithField("password", "testPassword"),
	})

	responseBodyRaw, errBody := ioutil.ReadAll(response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func TestDefaultApiClient_SendRequestWithPostAndFileUploadBody(t *testing.T) {
	testFileName := "test.Unit"
	testFileContent := "someFileContent"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotNilf(t, r.Body, "body should not be nil")
		pError := r.ParseMultipartForm(0)
		assert.Nilf(t, pError, "parsing form should not contain errors")

		file, err := DownloadFile(r, "file")

		assert.Nilf(t, err, "error should be null")
		assert.Equalf(t, testFileName, file.Name, "want %v found %v", testFileName, file.Name)

		assert.Equalf(t, testFileContent, string(file.File), "wrong file content, want %v, found %v", testFileContent, string(file.File))

		assert.Equalf(t, "POST", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")

		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultApiClient{
		Client: server.Client(),
	}

	if helper.FileExists(testFileName) {
		helper.DeleteFile(testFileName)
	}

	helper.WriteToFile(testFileContent, testFileName)

	response, err := api.SendRequest(ApiClientOptions{
		Method:   POST,
		Protocol: "http",
		Host:     server.URL,
		Path:     "/foo/bar",
		Body:     *NewApiClientBody().FormData().WithFile("file", testFileName),
	})

	helper.DeleteFile(testFileName)

	responseBodyRaw, errBody := ioutil.ReadAll(response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func TestDefaultApiClient_SendRequestWithPostAndFileUploadAndFieldsBody(t *testing.T) {
	testFileName := "test.Unit"
	testFileContent := "someFileContent"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotNilf(t, r.Body, "body should not be nil")
		var body TestBody
		pError := r.ParseMultipartForm(0)
		assert.Nilf(t, pError, "parsing form should not contain errors")

		body.Password = r.Form.Get("password")
		body.Username = r.Form.Get("username")

		file, err := DownloadFile(r, "file")

		assert.Nilf(t, err, "error should be null")
		assert.Equalf(t, testFileName, file.Name, "want %v found %v", testFileName, file.Name)

		assert.Equalf(t, testFileContent, string(file.File), "wrong file content, want %v, found %v", testFileContent, string(file.File))
		assert.Equalf(t, "testUser", body.Username, "username = %v, want %v", body.Username, "testUser")
		assert.Equalf(t, "testPassword", body.Password, "password = %v, want %v", body.Password, "testPassword")

		assert.Equalf(t, "POST", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")

		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultApiClient{
		Client: server.Client(),
	}

	if helper.FileExists(testFileName) {
		helper.DeleteFile(testFileName)
	}

	helper.WriteToFile(testFileContent, testFileName)

	response, err := api.SendRequest(ApiClientOptions{
		Method:   POST,
		Protocol: "http",
		Host:     server.URL,
		Path:     "/foo/bar",
		Body:     *NewApiClientBody().FormData().WithFile("file", testFileName).WithField("username", "testUser").WithField("password", "testPassword"),
	})

	helper.DeleteFile(testFileName)

	responseBodyRaw, errBody := ioutil.ReadAll(response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}
