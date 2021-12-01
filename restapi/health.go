package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
)

// Login Generate a token for a valid user
func (c *HttpListener) Probe() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "Healthy"
		json.NewEncoder(w).Encode(response)
	}
}
