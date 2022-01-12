package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/identity/jwt"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/gorilla/mux"
)

// Introspection Validates a token for a user
func (c *AuthorizationControllers) Introspection() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := execution_context.Get()
		vars := mux.Vars(r)
		tenantId := vars["tenantId"]

		// if no tenant is set we will assume it is the global tenant
		if tenantId == "" {
			tenantId = "global"
		}

		ctx.Authorization.TenantId = tenantId

		token := r.FormValue("token")

		if token == "" {
			response := models.OAuthErrorResponse{
				Error:            models.OAuthInvalidRequestError,
				ErrorDescription: "The JWT token was not found or the header was malformed, please check your authorization header",
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			logger.Error("There was an error validating token")
			return
		}

		userToken, err := jwt.ValidateUserToken(token, ctx.Authorization.Scope, ctx.Authorization.Audiences...)

		if err != nil {
			response := models.OAuthIntrospectResponse{
				Active: false,
			}

			json.NewEncoder(w).Encode(response)
			logger.Error("Token for user %v is not valid, %v", userToken.DisplayName, err.Error())
			return
		}

		response := models.OAuthIntrospectResponse{
			Active:    true,
			ID:        userToken.ID,
			TokenType: userToken.Scope,
			Subject:   userToken.User,
			ExpiresAt: fmt.Sprintf("%v", userToken.ExpiresAt.Unix()),
			IssuedAt:  fmt.Sprintf("%v", userToken.IssuedAt.Unix()),
			Issuer:    userToken.Issuer,
		}

		logger.Success("Token for user %v was validated successfully", userToken.DisplayName)
		json.NewEncoder(w).Encode(response)
	}
}
