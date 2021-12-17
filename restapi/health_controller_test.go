package restapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthController(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/probe", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(globalHttpListener.Probe())

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, "\"Healthy\"\n", rr.Body.String())
}
