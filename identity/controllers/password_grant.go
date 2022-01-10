package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/identity/jwt"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/security"
)

type PasswordGrantFlow struct{}

func (passwordGrantFlow PasswordGrantFlow) Authenticate(controller *AuthorizationControllers, request *models.OAuthLoginRequest, tenant string) (*models.OAuthLoginResponse, *models.OAuthErrorResponse) {
	var errorResponse models.OAuthErrorResponse
	ctx := execution_context.Get()
	user := controller.UserAdapter.GetUserByEmail(request.Username)

	if user.ID == "" {
		if user.Email == "" {
			errorResponse = models.OAuthErrorResponse{
				Error:            models.OAuthInvalidClientError,
				ErrorDescription: fmt.Sprintf("User %v was not found", user.DisplayName),
			}
		} else {
			errorResponse = models.OAuthErrorResponse{
				Error:            models.OAuthInvalidClientError,
				ErrorDescription: "Unknown error",
			}
		}
		logger.Error(errorResponse.ErrorDescription)
		return nil, &errorResponse
	}

	password := security.SHA256Encode(request.Password)

	if password != user.Password {
		errorResponse = models.OAuthErrorResponse{
			Error:            models.OAuthInvalidClientError,
			ErrorDescription: fmt.Sprintf("Invalid password for user %v", user.DisplayName),
		}
		logger.Error(errorResponse.ErrorDescription)
		return nil, &errorResponse
	}

	token, err := jwt.GenerateDefaultUserToken(*user)
	if err != nil {
		errorResponse = models.OAuthErrorResponse{
			Error:            models.OAuthInvalidClientError,
			ErrorDescription: fmt.Sprintf("There was an error validating user token, %v", err.Error()),
		}
		return nil, &errorResponse
	}

	controller.UserAdapter.UpdateUserRefreshToken(user.ID, token.RefreshToken)

	expiresIn := ctx.Authorization.Options.TokenDuration * 60
	response := models.OAuthLoginResponse{
		AccessToken:  token.Token,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    fmt.Sprintf("%v", expiresIn),
		TokenType:    "Bearer",
		Scope:        ctx.Authorization.Options.Scope,
	}

	logger.Success("Token for user %v was generated successfully", user.Username)

	return &response, nil
}

func (passwordGrantFlow PasswordGrantFlow) RefreshToken(controller *AuthorizationControllers, request *models.OAuthLoginRequest, tenant string) (*models.OAuthLoginResponse, *models.OAuthErrorResponse) {
	var errorResponse models.OAuthErrorResponse
	ctx := execution_context.Get()
	userEmail := jwt.GetTokenClaim(request.RefreshToken, "sub")
	user := controller.UserAdapter.GetUserByEmail(userEmail)

	if user.ID == "" {
		if user.DisplayName != "" {
			errorResponse = models.OAuthErrorResponse{
				Error:            models.OAuthInvalidClientError,
				ErrorDescription: fmt.Sprintf("User %v was not found", user.DisplayName),
			}
		} else {
			errorResponse = models.OAuthErrorResponse{
				Error:            models.OAuthInvalidClientError,
				ErrorDescription: "Unknown error",
			}
		}
		logger.Error(errorResponse.ErrorDescription)
		return nil, &errorResponse
	}

	userRefreshToken := user.RefreshToken
	if !strings.EqualFold(request.RefreshToken, userRefreshToken) {
		errorResponse = models.OAuthErrorResponse{
			Error:            models.OAuthInvalidClientError,
			ErrorDescription: "Refresh token is invalid",
		}
		logger.Error(errorResponse.ErrorDescription)
		return nil, &errorResponse
	}

	token, err := jwt.ValidateRefreshToken(request.RefreshToken, user.Email)
	if err != nil {
		errorResponse = models.OAuthErrorResponse{
			Error:            models.OAuthInvalidClientError,
			ErrorDescription: fmt.Sprintf("There was an error validating user token, %v", err.Error()),
		}
		return nil, &errorResponse
	}

	newToken, err := jwt.GenerateDefaultUserToken(*user)
	if err != nil {
		errorResponse = models.OAuthErrorResponse{
			Error:            models.OAuthInvalidClientError,
			ErrorDescription: fmt.Sprintf("There was an error generating the new user token, %v", err.Error()),
		}
		return nil, &errorResponse
	}

	expiresIn := ctx.Authorization.Options.TokenDuration * 60
	response := models.OAuthLoginResponse{
		AccessToken:  newToken.Token,
		RefreshToken: request.RefreshToken,
		ExpiresIn:    fmt.Sprintf("%v", expiresIn),
		TokenType:    "Bearer",
		Scope:        ctx.Authorization.Options.Scope,
	}

	todayPlus30 := time.Now().Add((time.Hour * 24) * 30)
	if token.ExpiresAt.Before(todayPlus30) {
		response.RefreshToken = newToken.RefreshToken
		controller.UserAdapter.UpdateUserRefreshToken(user.ID, newToken.RefreshToken)
	}

	logger.Success("Token for user %v was generated successfully", user.Username)

	return &response, nil
}
