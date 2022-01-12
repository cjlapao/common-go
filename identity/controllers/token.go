package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/service_provider"
	"github.com/gorilla/mux"
)

// Login Generate a token for a valid user
func (c *AuthorizationControllers) Token() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var errorResponse models.OAuthErrorResponse
		ctx := execution_context.Get()
		vars := mux.Vars(r)
		tenantId := vars["tenantId"]

		var loginRequest models.OAuthLoginRequest
		http_helper.MapRequestBody(r, &loginRequest)

		// if no tenant is set we will assume it is the global tenant
		if tenantId == "" {
			tenantId = "global"
		}

		baseUrl := service_provider.Get().GetBaseUrl(r)
		ctx.Authorization.TenantId = tenantId

		ctx.Authorization.Issuer = baseUrl + "/auth/" + tenantId

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

		switch loginRequest.GrantType {
		case "password":
			response, errorResponse := PasswordGrantFlow{}.Authenticate(c, &loginRequest, ctx.Authorization.TenantId)
			if errorResponse != nil {
				switch errorResponse.Error {
				case models.OAuthInvalidClientError:
					w.WriteHeader(http.StatusUnauthorized)
				default:
					w.WriteHeader(http.StatusBadRequest)
				}
				json.NewEncoder(w).Encode(*errorResponse)
				return
			}
			json.NewEncoder(w).Encode(*response)
			return
		case "refresh_token":
			response, errorResponse := PasswordGrantFlow{}.RefreshToken(c, &loginRequest, ctx.Authorization.TenantId)
			if errorResponse != nil {
				switch errorResponse.Error {
				case models.OAuthInvalidClientError:
					w.WriteHeader(http.StatusUnauthorized)
				default:
					w.WriteHeader(http.StatusBadRequest)
				}
				json.NewEncoder(w).Encode(*errorResponse)
				return
			}
			json.NewEncoder(w).Encode(*response)
			return
		default:
			w.WriteHeader(http.StatusBadRequest)
			errorResponse = models.OAuthErrorResponse{
				Error:            models.OAuthUnsupportedGrantType,
				ErrorDescription: fmt.Sprintf("Grant %v is not currently supported by the system", loginRequest.GrantType),
			}
			logger.Error(errorResponse.ErrorDescription)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}
}
