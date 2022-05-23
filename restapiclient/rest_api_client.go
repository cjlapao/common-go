package restapiclient

import (
	"context"
	"net/url"
)

type RestApiClient interface {
	Get(url string) RestApiClient
	Post(url string, body RestApiClientBody) RestApiClient
	PostForm(url string, values url.Values) RestApiClient
	Put(url string, body RestApiClientBody) RestApiClient
	Delete(url string) RestApiClient
	DeleteWithBody(url string, body RestApiClientBody) RestApiClient

	PreFlight(requestURL string, origins string, methods []RestApiClientMethod, headers []string) (*RestApiClientResponse, error)

	AddAuthorization(auth *RestApiClientAuthorization) RestApiClient
	AddBearerToken(token string) RestApiClient
	AddApiKey(key string) RestApiClient
	AddBasicAuth(username string, password string) RestApiClient

	Run(ctx ...context.Context) (*RestApiClientResponse, error)

	SendRequest(options RestApiClientRequest) (*RestApiClientResponse, error)
	SendRequestWithContext(apiRequest RestApiClientRequest, ctx context.Context) (*RestApiClientResponse, error)
}
