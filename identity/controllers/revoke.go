package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/common-go/identity/models"
)

// Revoke Revokes a user or a user refresh tenant, when revoking a user
// this will remove the user from the database deleting it
func (c *AuthorizationControllers) Revoke() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var revokeRequest models.OAuthRevokeRequest

		http_helper.MapRequestBody(r, &revokeRequest)

		if revokeRequest.ClientID == "" {
			w.WriteHeader(http.StatusBadRequest)
			ErrEmptyUserID.Log()
			json.NewEncoder(w).Encode(ErrEmptyUserID)
			return
		}

		dbUser := c.Context.UserDatabaseAdapter.GetUserById(revokeRequest.ClientID)

		if dbUser == nil || dbUser.ID == "" {
			w.WriteHeader(http.StatusBadRequest)
			ErrUserNotFound.Log()
			json.NewEncoder(w).Encode(ErrUserNotFound)
			return
		}

		// There is no revoke token filled in so we will be removing the user
		// otherwise we will only revoke the refresh token
		switch revokeRequest.GrantType {
		case "revoke_user":
			removeResult := c.Context.UserDatabaseAdapter.RemoveUser(dbUser.ID)

			if !removeResult {
				w.WriteHeader(http.StatusBadRequest)
				ErrUserNotRemoved.Log()
				json.NewEncoder(w).Encode(ErrUserNotRemoved)
				return
			}
		case "revoke_token":
			c.Context.UserDatabaseAdapter.UpdateUserRefreshToken(dbUser.ID, "")
		default:
			w.WriteHeader(http.StatusBadRequest)
			ErrGrantNotSupported.Log()
			json.NewEncoder(w).Encode(ErrGrantNotSupported)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
