package restapiclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RestApiClientResponse_Status(t *testing.T) {
	response := http.Response{
		Status:     "OK",
		StatusCode: 200,
	}

	apiResponse := RestApiClientResponse{
		Response: &response,
	}

	assert.Equal(t, "OK", apiResponse.Status())
}

func Test_RestApiClientResponse_StatusCode(t *testing.T) {
	response := http.Response{
		Status:     "OK",
		StatusCode: 200,
	}

	apiResponse := RestApiClientResponse{
		Response: &response,
	}

	assert.Equal(t, 200, apiResponse.StatusCode())
}

func Test_RestApiClientResponse_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		want     bool
	}{
		{
			name: "100 response",
			response: &http.Response{
				Status:     "OK",
				StatusCode: 101,
			},
			want: false,
		},
		{
			name: "200 response",
			response: &http.Response{
				Status:     "OK",
				StatusCode: 200,
			},
			want: true,
		},
		{
			name: "300 response",
			response: &http.Response{
				Status:     "OK",
				StatusCode: 307,
			},
			want: false,
		},
		{
			name: "400 responses",
			response: &http.Response{
				Status:     "OK",
				StatusCode: 404,
			},
			want: false,
		},
		{
			name: "500 response",
			response: &http.Response{
				Status:     "OK",
				StatusCode: 500,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiResponse := RestApiClientResponse{
				Response: tt.response,
			}

			assert.Equal(t, tt.want, apiResponse.IsSuccessful())
		})
	}
}
