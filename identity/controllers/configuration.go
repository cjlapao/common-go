package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/service_provider"
	"github.com/gorilla/mux"
)

// Login Generate a token for a valid user
func (c *AuthorizationControllers) Configuration() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tenantId := vars["tenantId"]

		// if no tenant is set we will assume it is the global tenant
		if tenantId == "" {
			tenantId = "global"
		}

		baseurl := service_provider.Get().GetBaseUrl(r)

		response := models.OAuthConfigurationResponse{
			Issuer:                baseurl + "/auth/" + tenantId,
			JwksURI:               baseurl + "/auth/" + tenantId + "/auth/.well-known/openid-configuration/jwks",
			AuthorizationEndpoint: baseurl + "/auth/" + tenantId + "/auth/authorize",
			TokenEndpoint:         baseurl + "/auth/" + tenantId + "/auth/token",
			UserinfoEndpoint:      baseurl + "/auth/" + tenantId + "/auth/userinfo",
			IntrospectionEndpoint: baseurl + "/auth/" + tenantId + "/auth/introspection",
		}

		json.NewEncoder(w).Encode(response)
	}
}
