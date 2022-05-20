package restapiclient

import (
	"context"
	"net/url"
)

type RestApiClient interface {
	Get(url string) RestApiClient
	Post(url string, body RestApiClientBody) RestApiClient
	PostForm(url string, values url.Values) RestApiClient
	PreFlight(requestURL string, methods []RestApiClientMethod, origins string, headers []string) (*RestApiClientResponse, error)

	AddAuthorization(auth *RestApiClientAuthorization) RestApiClient
	AddBearerToken(token string) RestApiClient
	AddApiKey(key string) RestApiClient
	AddBasicAuth(username string, password string) RestApiClient

	Run(ctx ...context.Context) (*RestApiClientResponse, error)

	SendRequest(options RestApiClientRequest) (*RestApiClientResponse, error)
	SendRequestWithContext(apiRequest RestApiClientRequest, ctx context.Context) (*RestApiClientResponse, error)
}
