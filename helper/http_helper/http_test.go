package http_helper

import (
	"net/http"
	"testing"
)

func TestGetContentType(t *testing.T) {
	tests := []struct {
		name string
		args http.Header
		want string
	}{
		{
			"content-type only",
			http.Header{
				"Content-Type": {"application/json", "something else"},
				"Other-Header": {"someKey", "someKey2"},
			},
			"application/json",
		},
		{
			"content-type only",
			http.Header{
				"Content-Type": {"application/json; charset"},
				"Other-Header": {"someKey", "someKey2"},
			},
			"application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetContentType(tt.args); got != tt.want {
				t.Errorf("GetContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAuthorizationToken(t *testing.T) {
	tests := []struct {
		name  string
		args  http.Header
		want  string
		want1 bool
	}{
		{
			"tes",
			http.Header{
				"authorization": {"Bearer 12"},
			},
			"12",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetAuthorizationToken(tt.args)
			if got != tt.want {
				t.Errorf("GetAuthorizationToken() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetAuthorizationToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
