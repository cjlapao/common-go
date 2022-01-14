package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/identity/oauthflow"
)

// Login Generate a token for a valid user
func (c *AuthorizationControllers) Token() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest models.OAuthLoginRequest
		http_helper.MapRequestBody(r, &loginRequest)

		switch loginRequest.GrantType {
		case "password":
			response, errorResponse := oauthflow.PasswordGrantFlow{}.Authenticate(&loginRequest)
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
			if loginRequest.Username != "" {
				response, errorResponse := oauthflow.PasswordGrantFlow{}.RefreshToken(&loginRequest)
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
			} else if loginRequest.ClientID != "" {
				// TODO: Implement client id validations
				w.WriteHeader(http.StatusBadRequest)
				ErrGrantNotSupported.Log()
				json.NewEncoder(w).Encode(ErrGrantNotSupported)
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				ErrGrantNotSupported.Log()
				json.NewEncoder(w).Encode(ErrGrantNotSupported)
				return

			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			ErrGrantNotSupported.Log()
			json.NewEncoder(w).Encode(ErrGrantNotSupported)
			return
		}
	}
}
