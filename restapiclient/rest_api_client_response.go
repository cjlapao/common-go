package restapiclient

import "net/http"

type RestApiClientResponse struct {
	Response *http.Response // http raw response

}

// Status Gets the response status as a string
func (r RestApiClientResponse) Status() string {
	if r.Response == nil {
		return ""
	}
	return r.Response.Status
}

// StatusCode Get's the response status code
func (r RestApiClientResponse) StatusCode() int {
	if r.Response == nil {
		return 0
	}

	return r.Response.StatusCode
}

// IsSuccessful checks if the response contains a successful status code
// this can range from 200 to 299. while 300-309 are not error codes, the call
// itself is not a successful as it will redirect to another call that needs validation
func (r RestApiClientResponse) IsSuccessful() bool {
	if r.StatusCode() >= 200 && r.StatusCode() <= 299 {
		return true
	}

	return false
}
