package restapi

import (
	"encoding/json"
	"net/http"
)

// Login Generate a token for a valid user
func (c *HttpListener) Probe(w http.ResponseWriter, r *http.Request) {
	response := "Healthy"
	json.NewEncoder(w).Encode(response)
}
