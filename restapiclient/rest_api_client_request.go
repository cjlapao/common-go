package restapiclient

import (
	"errors"
	"net/url"
)

type RestApiClientRequest struct {
	Method        RestApiClientMethod
	Authorization *RestApiClientAuthorization
	URL           *url.URL
	Body          *RestApiClientBody
	Headers       map[string]string
}

func NewRestApiRequest(url string, method RestApiClientMethod) *RestApiClientRequest {
	request := RestApiClientRequest{
		Method:  method,
		Headers: map[string]string{},
	}

	err := request.ParseUrl(url)

	if err != nil {
		return nil
	}

	return &request
}

func (request *RestApiClientRequest) AddHeader(key string, value string) error {
	if request.Headers == nil {
		request.Headers = map[string]string{}
	}

	if key == "" {
		return errors.New("key cannot be empty")
	}

	if value == "" {
		return errors.New("value cannot be empty")
	}

	request.Headers[key] = value
	return nil
}

func (request *RestApiClientRequest) ParseUrl(value string) error {
	var parsedURL *url.URL
	var err error

	if parsedURL, err = url.Parse(value); err != nil {
		return err
	}

	request.URL = parsedURL
	return nil
}
