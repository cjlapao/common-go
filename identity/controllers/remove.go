package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/security"
	"github.com/cjlapao/common-go/service_provider"
	"github.com/gorilla/mux"
)

// Register Create an user in the tenant
func (c *AuthorizationControllers) Remove() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var errorResponse models.OAuthErrorResponse
		var registerRequest models.OAuthRegisterRequest
		ctx := execution_context.Get()
		vars := mux.Vars(r)
		tenantId := vars["tenantId"]

		if c.UserAdapter == nil {
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse = models.OAuthErrorResponse{
				Error:            models.OAuthInvalidRequestError,
				ErrorDescription: "Context is null and cannot be validated",
			}
			logger.Error(errorResponse.ErrorDescription)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		http_helper.MapRequestBody(r, &registerRequest)

		// if no tenant is set we will assume it is the global tenant
		if tenantId == "" {
			tenantId = "global"
		}

		baseUrl := service_provider.Get().GetBaseUrl(r)
		ctx.Authorization.TenantId = tenantId
		ctx.Authorization.Options.Issuer = baseUrl + "/auth/" + tenantId

		user := models.NewUser()
		user.Username = registerRequest.Username
		user.Email = registerRequest.Email
		user.FirstName = registerRequest.FirstName
		user.LastName = registerRequest.LastName
		user.DisplayName = user.FirstName + " " + user.LastName
		user.Password = security.SHA256Encode(registerRequest.Password)
		user.InvalidAttempts = 0
		user.EmailVerified = false
		if registerRequest.Claims != nil && len(registerRequest.Claims) > 0 {
			for _, claim := range registerRequest.Claims {
				user.Claims = append(user.Claims, models.NewUserClaim(claim))
			}
		} else {
			user.Claims = append(user.Claims, models.ReadClaim)
		}

		if registerRequest.Roles != nil && len(registerRequest.Roles) > 0 {
			for _, role := range registerRequest.Roles {
				user.Roles = append(user.Roles, models.NewUserRole(role))
			}
		} else {
			user.Roles = append(user.Roles, models.RegularUserRole)
		}

		if !user.IsValid() {
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse = models.OAuthErrorResponse{
				Error:            models.OAuthInvalidClientError,
				ErrorDescription: "User is not valid",
			}
			logger.Error(errorResponse.ErrorDescription)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		dbUser := c.UserAdapter.GetUserByEmail(user.Email)
		if dbUser.Email != "" {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse = models.OAuthErrorResponse{
				Error:            models.OAuthInvalidRequestError,
				ErrorDescription: fmt.Sprintf("User %v already exists on tenant %v", user.Email, tenantId),
			}
			logger.Error(errorResponse.ErrorDescription)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		c.UserAdapter.UpsertUser(*user)

		json.NewEncoder(w).Encode(user)
	}
}
