package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/automapper"
	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/gorilla/mux"
)

// Login Generate a token for a valid user
func (c *AuthorizationControllers) Jwks() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := execution_context.Get()
		vars := mux.Vars(r)
		tenantId := vars["tenantId"]

		// if no tenant is set we will assume it is the global tenant
		if tenantId == "" {
			tenantId = "global"
		}

		response := models.OAuthJwksResponse{
			Keys: make([]models.OAuthJwksKey, 0),
		}
		key := models.OAuthJwksKey{}
		defaultKey := ctx.Authorization.KeyVault.GetDefaultKey()
		if len(defaultKey.JWK.Keys) >= 1 {
			automapper.Map(defaultKey.JWK.Keys[0], &key)
		}

		key.ID = defaultKey.ID
		response.Keys = append(response.Keys, key)

		json.NewEncoder(w).Encode(response)
	}
}
