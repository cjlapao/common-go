package apiclient

import "net/http"

type ApiClient interface {
	GET(host string, path string, token string) (interface{}, error)
	POST(host string, path string, token string) (interface{}, error)
}

type DefaultApiClient struct {
	Client *http.Client
}

func NewApiClient() ApiClient {
	client := &http.Client{}
	apiClient := DefaultApiClient{
		Client: client,
	}

	return &apiClient
}

func (c DefaultApiClient) GET(host string, path string, token string) (interface{}, error) {
	return nil, nil
}

func (c DefaultApiClient) POST(host string, path string, token string) (interface{}, error) {
	return nil, nil
}

func (c *DefaultApiClient) sendRequest(method string, uri string, token string, body interface{}) {

}
