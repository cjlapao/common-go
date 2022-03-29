package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/service_provider"
)

// Configuration Returns the OpenID Oauth configuration endpoint
func (c *AuthorizationControllers) Configuration() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		baseurl := service_provider.Get().GetBaseUrl(r)

		response := models.OAuthConfigurationResponse{
			Issuer:                baseurl + "/auth/" + c.Context.Authorization.TenantId,
			JwksURI:               baseurl + "/auth/" + c.Context.Authorization.TenantId + "/auth/.well-known/openid-configuration/jwks",
			AuthorizationEndpoint: baseurl + "/auth/" + c.Context.Authorization.TenantId + "/auth/authorize",
			TokenEndpoint:         baseurl + "/auth/" + c.Context.Authorization.TenantId + "/auth/token",
			UserinfoEndpoint:      baseurl + "/auth/" + c.Context.Authorization.TenantId + "/auth/userinfo",
			IntrospectionEndpoint: baseurl + "/auth/" + c.Context.Authorization.TenantId + "/auth/introspection",
		}

		json.NewEncoder(w).Encode(response)
	}
}
