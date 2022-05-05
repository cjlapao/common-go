package apiclient

import (
	"net/http"

	"github.com/cjlapao/common-go/guard"
)

type ApiClient interface {
	GET(host string, path string, options ...ApiClientOptions) (interface{}, error)
	POST(host string, path string, options ...ApiClientOptions) (interface{}, error)

	SendRequest(options ApiClientOptions) (*http.Response, error)
}

type DefaultApiClient struct {
	Client    *http.Client
	Transport *http.Transport
}

func NewApiClient() ApiClient {
	client := &http.Client{}
	apiClient := DefaultApiClient{
		Client: client,
	}

	return &apiClient
}

func (c DefaultApiClient) GET(host string, path string, options ...ApiClientOptions) (interface{}, error) {
	return nil, nil
}

func (c DefaultApiClient) POST(host string, path string, options ...ApiClientOptions) (interface{}, error) {
	return nil, nil
}

func (c *DefaultApiClient) SendRequest(options ApiClientOptions) (*http.Response, error) {
	if err := guard.EmptyOrNil(options.Method.String()); err != nil {
		options.Method = GET
	}

	url, err := options.Url()
	if err != nil {
		return nil, err
	}

	var request *http.Request

	if guard.EmptyOrNil(options.Body); err != nil {
		request, err = http.NewRequest(options.Method.String(), url.String(), nil)
	} else {
		request = &http.Request{
			Method: options.Method.String(),
			URL:    url,
		}
	}

	response, err := c.Client.Do(request)

	if err != nil {
		return nil, err
	}

	return response, nil
}
