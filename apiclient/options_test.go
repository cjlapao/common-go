package apiclient

import (
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_URL(t *testing.T) {
	// arrange
	expectedUrl, _ := url.Parse("http://example.com")
	expectedUrlWithPath, _ := url.Parse("http://example.com/foobar")
	expectedUrlWithOneQuery, _ := url.Parse("http://example.com/foobar?foo=bar")
	expectedUrlWithQueries, _ := url.Parse("http://example.com/foobar?foo=bar&bar=foo")
	expectedUrlWithQueriesButNoPath, _ := url.Parse("http://example.com?foo=bar&bar=foo&a=o")

	tests := []struct {
		name      string
		options   ApiClientOptions
		want      *url.URL
		expectErr bool
	}{
		{
			"with invalid url",
			ApiClientOptions{
				Method:   GET,
				Host:     "",
				Protocol: "http",
				Path:     "/",
			},
			expectedUrl,
			true,
		},
		{
			"with invalid scheme",
			ApiClientOptions{
				Method:   GET,
				Host:     "",
				Protocol: "httpwhatnot",
				Path:     "/",
			},
			expectedUrl,
			true,
		},
		{
			"with trailing /",
			ApiClientOptions{
				Method:   GET,
				Host:     "example.com",
				Protocol: "http",
				Path:     "/",
			},
			expectedUrl,
			false,
		},
		{
			"with fragment",
			ApiClientOptions{
				Method:   GET,
				Host:     "example.com",
				Protocol: "http",
				Path:     "/foobar",
			},
			expectedUrlWithPath,
			false,
		},
		{
			"with query parameter and fragment",
			ApiClientOptions{
				Method:   GET,
				Host:     "example.com",
				Protocol: "http",
				Path:     "/foobar",
				QueryParameters: map[string]string{
					"foo": "bar",
				},
			},
			expectedUrlWithOneQuery,
			false,
		},
		{
			"with query parameters and multiple fragments",
			ApiClientOptions{
				Method:   GET,
				Host:     "example.com",
				Protocol: "http",
				Path:     "/foobar",
				QueryParameters: map[string]string{
					"foo": "bar",
					"bar": "foo",
				},
			},
			expectedUrlWithQueries,
			false,
		},
		{
			"with query parameters but no fragment",
			ApiClientOptions{
				Method:   GET,
				Host:     "example.com",
				Protocol: "http",
				QueryParameters: map[string]string{
					"foo": "bar",
					"bar": "foo",
					"a":   "o",
				},
			},
			expectedUrlWithQueriesButNoPath,
			false,
		},
		{
			"with complex host and https",
			ApiClientOptions{
				Method: GET,
				Host:   "http://example.com",
			},
			&url.URL{
				Host:   "example.com",
				Scheme: "http",
			},
			false,
		},
		{
			"with complex host, port and http",
			ApiClientOptions{
				Method: GET,
				Host:   "http://example.com:400",
			},
			&url.URL{
				Host:   "example.com:400",
				Scheme: "http",
			},
			false,
		},
		{
			"with complex host, port and https",
			ApiClientOptions{
				Method: GET,
				Host:   "https://example.com:400",
				QueryParameters: map[string]string{
					"foo": "bar",
				},
			},
			&url.URL{
				Host:     "example.com:400",
				Scheme:   "https",
				RawQuery: "foo=bar",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.options.Url()
			if (err != nil) != tt.expectErr {
				t.Errorf("ApiClientOptions.Url() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if (err != nil) && tt.expectErr {
				return
			}

			if !reflect.DeepEqual(got.Host, tt.want.Host) || !reflect.DeepEqual(got.Scheme, tt.want.Scheme) {
				t.Errorf("ApiClientOptions.Url() = %v, want %v", got, tt.want)
			}

			gotParts := strings.Split(got.RawQuery, "&")
			wantParts := strings.Split(tt.want.RawQuery, "&")

			assert.Equal(t, len(wantParts), len(gotParts))
			for _, wantPart := range wantParts {
				found := false
				for _, gotPart := range gotParts {
					if wantPart == gotPart {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ApiClientOptions.Url().Part = %v, not found", wantPart)
				}
			}
		})
	}
}
