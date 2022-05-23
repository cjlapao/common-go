package restapiclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cjlapao/common-go/guard"
)

type DefaultRestApiClient struct {
	Client    *http.Client
	Transport *http.Transport
	request   *RestApiClientRequest
}

func NewRestApiClient() RestApiClient {
	client := &http.Client{}
	RestApiClient := DefaultRestApiClient{
		Client: client,
	}

	return &RestApiClient
}

func (c *DefaultRestApiClient) Get(requestURL string) RestApiClient {
	parsedURL, err := url.Parse(requestURL)

	if err != nil {
		return nil
	}

	c.request = &RestApiClientRequest{
		Method: API_METHOD_GET,
		URL:    parsedURL,
	}

	return c
}

func (c *DefaultRestApiClient) Post(requestURL string, body RestApiClientBody) RestApiClient {
	parsedURL, err := url.Parse(requestURL)

	if err != nil {
		return nil
	}

	c.request = &RestApiClientRequest{
		Method: API_METHOD_POST,
		URL:    parsedURL,
	}

	c.request.Body = &body
	return c
}

func (c *DefaultRestApiClient) Put(requestURL string, body RestApiClientBody) RestApiClient {
	parsedURL, err := url.Parse(requestURL)

	if err != nil {
		return nil
	}

	c.request = &RestApiClientRequest{
		Method: API_METHOD_PUT,
		URL:    parsedURL,
	}

	c.request.Body = &body
	return c
}

func (c *DefaultRestApiClient) Delete(requestURL string) RestApiClient {
	parsedURL, err := url.Parse(requestURL)

	if err != nil {
		return nil
	}

	c.request = &RestApiClientRequest{
		Method: API_METHOD_DELETE,
		URL:    parsedURL,
	}

	return c
}

func (c *DefaultRestApiClient) DeleteWithBody(requestURL string, body RestApiClientBody) RestApiClient {
	parsedURL, err := url.Parse(requestURL)

	if err != nil {
		return nil
	}

	c.request = &RestApiClientRequest{
		Method: API_METHOD_DELETE,
		URL:    parsedURL,
	}

	c.request.Body = &body
	return c
}

func (c *DefaultRestApiClient) PostForm(requestURL string, values url.Values) RestApiClient {
	parsedURL, err := url.Parse(requestURL)

	if err != nil {
		return nil
	}

	c.request = &RestApiClientRequest{
		Method: API_METHOD_POST,
		URL:    parsedURL,
	}

	c.request.Body = NewRestApiClientBody().UrlEncoded()
	for key, value := range values {
		c.request.Body.WithField(key, value[0])
	}

	return c
}

func (c *DefaultRestApiClient) PreFlight(requestURL string, origin string, methods []RestApiClientMethod, headers []string) (*RestApiClientResponse, error) {
	_, err := url.Parse(requestURL)

	if err != nil {
		return nil, err
	}

	if origin == "" {
		origin = "*"
	}

	c.request = NewRestApiRequest(requestURL, API_METHOD_OPTIONS)

	if c.request == nil {
		return nil, fmt.Errorf("error parsing url %v", requestURL)
	}

	parsedMethods := ""
	for _, method := range methods {
		if len(parsedMethods) > 0 {
			parsedMethods += ", "
		}
		parsedMethods += method.String()
	}

	if len(parsedMethods) > 0 {
		c.request.Headers["Access-Control-Request-Method"] = parsedMethods
	}

	parsedHeaders := ""
	for _, header := range headers {
		if len(parsedHeaders) > 0 {
			parsedHeaders += ", "
		}
		parsedHeaders += header
	}

	if len(parsedHeaders) > 0 {
		c.request.Headers["Access-Control-Request-Headers"] = parsedHeaders
	}

	c.request.Headers["Origin"] = origin

	response, err := c.Run()

	return response, err
}

func (c *DefaultRestApiClient) AddAuthorization(auth *RestApiClientAuthorization) RestApiClient {
	c.request.Authorization = auth

	return c
}

func (c *DefaultRestApiClient) AddBearerToken(token string) RestApiClient {
	c.request.Authorization = NewBearerTokenAuth(token)

	return c
}

func (c *DefaultRestApiClient) AddApiKey(key string) RestApiClient {
	c.request.Authorization = NewStandardApiKeyAuth(key)

	return c
}

func (c *DefaultRestApiClient) AddBasicAuth(username string, password string) RestApiClient {
	c.request.Authorization = NewBasicAuth(username, password)

	return c
}

func (c *DefaultRestApiClient) Run(ctx ...context.Context) (*RestApiClientResponse, error) {
	var apiResponse *RestApiClientResponse
	var err error

	if c.request == nil {
		return nil, errors.New("no request present")
	}

	if len(ctx) > 0 {
		apiResponse, err = c.SendRequestWithContext(*c.request, ctx[0])
	} else {
		apiResponse, err = c.SendRequest(*c.request)
	}

	return apiResponse, err
}

func (c *DefaultRestApiClient) SendRequest(options RestApiClientRequest) (*RestApiClientResponse, error) {
	return c.SendRequestWithContext(options, nil)
}

// SendRequestWithContext Sends an RestApiClientRequest and returns a
func (c *DefaultRestApiClient) SendRequestWithContext(apiRequest RestApiClientRequest, ctx context.Context) (*RestApiClientResponse, error) {
	var err error
	var response *http.Response

	if err = guard.EmptyOrNil(apiRequest.Method.String()); err != nil {
		apiRequest.Method = API_METHOD_GET
	}

	// Parsing the URL
	if apiRequest.URL == nil {
		return nil, errors.New("url cannot be nil")
	}

	var httpRequest *http.Request

	// Creating the request
	if apiRequest.Body == nil {
		httpRequest, err = http.NewRequest(apiRequest.Method.String(), apiRequest.URL.String(), nil)
		if err != nil {
			return nil, err
		}
	} else {
		httpRequest, err = http.NewRequest(apiRequest.Method.String(), apiRequest.URL.String(), apiRequest.Body.Get())
		if err != nil {
			return nil, err
		}

		key, value := apiRequest.Body.GetHeader()
		httpRequest.Header.Add(key, value)
	}

	// Adding all headers in the request
	for key, value := range apiRequest.Headers {
		httpRequest.Header.Add(key, value)
	}

	// Adding authorization token to the request
	if apiRequest.Authorization != nil {
		if apiRequest.Authorization.Value != "" {
			key, header := apiRequest.Authorization.GetHeader()
			httpRequest.Header[key] = header
		}
	}

	// Executing the call with context if present otherwise just a normal context
	if ctx == nil {
		response, err = c.Client.Do(httpRequest)
	} else {
		response, err = c.Client.Do(httpRequest.WithContext(ctx))
	}

	if err != nil {
		return nil, err
	}

	// creating the api response
	apiResponse := RestApiClientResponse{
		Response: response,
	}

	return &apiResponse, nil
}
