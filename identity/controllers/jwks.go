package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/automapper"
	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/identity/models"
)

// Jwks Returns the public keys for the openid oauth configuration endpoint for validation
func (c *AuthorizationControllers) Jwks() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		response := models.OAuthJwksResponse{
			Keys: make([]models.OAuthJwksKey, 0),
		}
		key := models.OAuthJwksKey{}
		defaultKey := c.Context.Authorization.KeyVault.GetDefaultKey()
		if len(defaultKey.JWK.Keys) >= 1 {
			automapper.Map(defaultKey.JWK.Keys[0], &key)
		}

		key.ID = defaultKey.ID
		response.Keys = append(response.Keys, key)

		json.NewEncoder(w).Encode(response)
	}
}
