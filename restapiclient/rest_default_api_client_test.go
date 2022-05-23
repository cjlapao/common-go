package restapiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/security"
	"github.com/stretchr/testify/assert"
)

type TestBody struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

func Test_DefaultRestApiClient_WithNoUrl_RaiseException(t *testing.T) {
	api := DefaultRestApiClient{}

	response, err := api.SendRequest(RestApiClientRequest{Method: API_METHOD_GET})

	assert.Nil(t, response)
	assert.NotNil(t, err)
}

func Test_DefaultRestApiClient_SendRequestWithGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "GET", r.Method, "Expected to be GET method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiRequest := RestApiClientRequest{
		Method: API_METHOD_GET,
	}

	apiRequest.ParseUrl(server.URL + "/foo/bar")

	apiResponse, err := api.SendRequest(apiRequest)

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_SendRequestWithGetAndTokenAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "GET", r.Method, "Expected to be GET method")
		token := r.Header.Get("Authorization")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		assert.NotNilf(t, token, "token should not be nil was %v", token)
		assert.Equalf(t, "Bearer abc", token, "Token should be \"Bearer abc\" found %v", token)
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiRequest := RestApiClientRequest{
		Method:        API_METHOD_GET,
		Authorization: NewBearerTokenAuth("abc"),
	}

	apiRequest.ParseUrl(server.URL + "/foo/bar")

	apiResponse, err := api.SendRequest(apiRequest)

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_SendRequestWithGetAndBasicAuthorization(t *testing.T) {
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

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiRequest := RestApiClientRequest{
		Method:        API_METHOD_GET,
		Authorization: NewBasicAuth("testUser", "testPassword"),
	}

	apiRequest.ParseUrl(server.URL + "/foo/bar")

	apiResponse, err := api.SendRequest(apiRequest)

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_SendRequestWithGetAndApiKeyAuthorization(t *testing.T) {
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

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiRequest := RestApiClientRequest{
		Method:        API_METHOD_GET,
		Authorization: NewApiKeyAuth("TestKey", "someKey"),
	}

	apiRequest.ParseUrl(server.URL + "/foo/bar")

	apiResponse, err := api.SendRequest(apiRequest)

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_SendRequestWithGetAndStandardApiKeyAuthorization(t *testing.T) {
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

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiRequest := RestApiClientRequest{
		Method:        API_METHOD_GET,
		Authorization: NewStandardApiKeyAuth("someKey"),
	}

	apiRequest.ParseUrl(server.URL + "/foo/bar")

	apiResponse, err := api.SendRequest(apiRequest)

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_SendRequestWithPostAndJsonBody(t *testing.T) {
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

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	testObj := TestBody{
		Username: "testUser",
		Password: "testPassword",
	}

	marshalledTestObj, _ := json.Marshal(testObj)

	apiRequest := RestApiClientRequest{
		Method: API_METHOD_POST,
		Body:   NewRestApiClientBody().Json(marshalledTestObj),
	}

	apiRequest.ParseUrl(server.URL + "/foo/bar")

	apiResponse, err := api.SendRequest(apiRequest)

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_SendRequestWithPostAndXUrlEncodedBody(t *testing.T) {
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

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiRequest := RestApiClientRequest{
		Method: API_METHOD_POST,
		Body:   NewRestApiClientBody().UrlEncoded().WithField("username", "testUser").WithField("password", "testPassword"),
	}

	apiRequest.ParseUrl(server.URL + "/foo/bar")

	apiResponse, err := api.SendRequest(apiRequest)

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_SendRequestWithPostAndFormDataBody(t *testing.T) {
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

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiRequest := RestApiClientRequest{
		Method: API_METHOD_POST,
		Body:   NewRestApiClientBody().FormData().WithField("username", "testUser").WithField("password", "testPassword"),
	}

	apiRequest.ParseUrl(server.URL + "/foo/bar")

	apiResponse, err := api.SendRequest(apiRequest)

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_SendRequestWithPostAndFileUploadBody(t *testing.T) {
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

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	if helper.FileExists(testFileName) {
		helper.DeleteFile(testFileName)
	}

	helper.WriteToFile(testFileContent, testFileName)

	apiRequest := RestApiClientRequest{
		Method: API_METHOD_POST,
		Body:   NewRestApiClientBody().FormData().WithFile("file", testFileName),
	}

	apiRequest.ParseUrl(server.URL + "/foo/bar")

	apiResponse, err := api.SendRequest(apiRequest)

	helper.DeleteFile(testFileName)

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_SendRequestWithPostAndFileUploadAndFieldsBody(t *testing.T) {
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

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	if helper.FileExists(testFileName) {
		helper.DeleteFile(testFileName)
	}

	helper.WriteToFile(testFileContent, testFileName)

	apiRequest := RestApiClientRequest{
		Method: API_METHOD_POST,
		Body:   NewRestApiClientBody().FormData().WithFile("file", testFileName).WithField("username", "testUser").WithField("password", "testPassword"),
	}

	apiRequest.ParseUrl(server.URL + "/foo/bar")

	apiResponse, err := api.SendRequest(apiRequest)

	helper.DeleteFile(testFileName)

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "GET", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")

		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiResponse, err := api.Get(server.URL + "/foo/bar").Run()

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_Get_WithTokenAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "GET", r.Method, "Expected to be GET method")
		token := r.Header.Get("Authorization")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		assert.NotNilf(t, token, "token should not be nil was %v", token)
		assert.Equalf(t, "Bearer abc", token, "Token should be \"Bearer abc\" found %v", token)

		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiResponse, err := api.Get(server.URL + "/foo/bar").
		AddBearerToken("abc").
		Run()

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_Get_WithApiKeyAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "GET", r.Method, "Expected to be GET method")
		apiKey := r.Header.Get("Authorization")
		expectedApiKey := "ApiKey abc"
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		assert.NotNilf(t, apiKey, "token should not be nil")
		assert.Equalf(t, expectedApiKey, apiKey, "Authorization = %v, want %v", apiKey, expectedApiKey)

		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiResponse, err := api.Get(server.URL + "/foo/bar").
		AddApiKey("abc").
		Run()

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_Get_WithBasicAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "GET", r.Method, "Expected to be GET method")
		token := r.Header.Get("Authorization")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")

		user, _ := security.EncodeString("testUser:testPassword")
		assert.Equalf(t, fmt.Sprintf("Basic %v", user), token, "Authorization = %v, want %v", token, fmt.Sprintf("Basic %v", user))
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiResponse, err := api.Get(server.URL+"/foo/bar").
		AddBasicAuth("testUser", "testPassword").
		Run()

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_Run_WithNoRequest(t *testing.T) {
	api := DefaultRestApiClient{}

	response, err := api.Run()

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "no request present", err.Error())
}

func Test_DefaultRestApiClient_Post_FormDataBody(t *testing.T) {
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

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	body := NewRestApiClientBody().FormData().WithField("username", "testUser").WithField("password", "testPassword")
	apiResponse, err := api.Post(server.URL+"/foo/bar", *body).
		Run()

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_Post_FileUpload(t *testing.T) {
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

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	if helper.FileExists(testFileName) {
		helper.DeleteFile(testFileName)
	}

	helper.WriteToFile(testFileContent, testFileName)

	apiResponse, err := api.Post(server.URL+"/foo/bar", *NewRestApiClientBody().FormData().WithFile("file", testFileName)).
		Run()

	helper.DeleteFile(testFileName)

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_PostForm(t *testing.T) {
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

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	form := url.Values{}
	form.Add("username", "testUser")
	form.Add("password", "testPassword")

	apiResponse, err := api.PostForm(server.URL+"/foo/bar", form).Run()

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_PreFlight(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equalf(t, "OPTIONS", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		assert.Equal(t, []string{"GET, POST"}, r.Header["Access-Control-Request-Method"])
		assert.Equal(t, []string{"X-Requested-With"}, r.Header["Access-Control-Request-Headers"])
		assert.Equal(t, []string{"*"}, r.Header["Origin"])

		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiResponse, err := api.PreFlight(server.URL+"/foo/bar", "*", []RestApiClientMethod{API_METHOD_GET, API_METHOD_POST}, []string{"X-Requested-With"})

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_PreFlight_WithOnlyMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equalf(t, "OPTIONS", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		assert.Equal(t, []string{"GET, POST"}, r.Header["Access-Control-Request-Method"])
		assert.Nil(t, r.Header["Access-Control-Request-Headers"])
		assert.Equal(t, []string{"*"}, r.Header["Origin"])

		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiResponse, err := api.PreFlight(server.URL+"/foo/bar", "*", []RestApiClientMethod{API_METHOD_GET, API_METHOD_POST}, []string{})

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_PreFlight_WithOnlyHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equalf(t, "OPTIONS", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")
		assert.Nil(t, r.Header["Access-Control-Request-Method"])
		assert.Equal(t, []string{"X-Requested-With"}, r.Header["Access-Control-Request-Headers"])
		assert.Equal(t, []string{"*"}, r.Header["Origin"])

		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiResponse, err := api.PreFlight(server.URL+"/foo/bar", "*", []RestApiClientMethod{}, []string{"X-Requested-With"})

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_Put(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotNilf(t, r.Body, "body should not be nil")
		var body TestBody
		pError := r.ParseMultipartForm(0)
		assert.Nilf(t, pError, "parsing form should not contain errors")

		body.Password = r.Form.Get("password")
		body.Username = r.Form.Get("username")

		assert.Equalf(t, "PUT", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")

		assert.Equalf(t, "testUser", body.Username, "username = %v, want %v", body.Username, "testUser")
		assert.Equalf(t, "testPassword", body.Password, "password = %v, want %v", body.Password, "testPassword")
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	body := NewRestApiClientBody().FormData().WithField("username", "testUser").WithField("password", "testPassword")
	apiResponse, err := api.Put(server.URL+"/foo/bar", *body).
		Run()

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equalf(t, "DELETE", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")

		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	apiResponse, err := api.Delete(server.URL + "/foo/bar").Run()

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}

func Test_DefaultRestApiClient_DeleteWithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotNilf(t, r.Body, "body should not be nil")
		var body TestBody
		pError := r.ParseMultipartForm(0)
		assert.Nilf(t, pError, "parsing form should not contain errors")

		body.Password = r.Form.Get("password")
		body.Username = r.Form.Get("username")

		assert.Equalf(t, "DELETE", r.Method, "Expected to be POST method")
		assert.Equalf(t, "/foo/bar", r.URL.String(), "Expected url to be /foo/bar")

		assert.Equalf(t, "testUser", body.Username, "username = %v, want %v", body.Username, "testUser")
		assert.Equalf(t, "testPassword", body.Password, "password = %v, want %v", body.Password, "testPassword")
		json.NewEncoder(w).Encode("ok")
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	body := NewRestApiClientBody().FormData().WithField("username", "testUser").WithField("password", "testPassword")
	apiResponse, err := api.DeleteWithBody(server.URL+"/foo/bar", *body).
		Run()

	responseBodyRaw, errBody := ioutil.ReadAll(apiResponse.Response.Body)
	responseBody := string(responseBodyRaw)
	assert.Nilf(t, err, "Response error should be nil")
	assert.Nilf(t, errBody, "Response body should not be nil")
	assert.Equalf(t, "\"ok\"\n", responseBody, "Response body should be ok")
}
