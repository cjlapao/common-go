package restapiclient

import (
	"net/url"
)

type RestApiClientRequest struct {
	Method        RestApiClientMethod
	Authorization *RestApiClientAuthorization
	URL           *url.URL
	Body          *RestApiClientBody
	Headers       map[string]string
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
