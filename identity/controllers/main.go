package controllers

import (
	"github.com/cjlapao/common-go/identity/database"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/log"
)

var logger = log.Get()

type AuthorizationFlow interface {
	Authorize(controller *AuthorizationControllers, request *models.OAuthLoginRequest, tenant string) (*models.OAuthLoginResponse, *models.OAuthErrorResponse)
	RefreshToken(controller *AuthorizationControllers, request *models.OAuthLoginRequest, tenant string) (*models.OAuthLoginResponse, *models.OAuthErrorResponse)
}

type AuthorizationControllers struct {
	UserAdapter database.UserDatabaseAdapter
}

func NewDefaultAuthorizationControllers() *AuthorizationControllers {
	context := database.NewMemoryUserAdapter()
	controllers := AuthorizationControllers{
		UserAdapter: context,
	}
	return &controllers
}

func NewAuthorizationControllers(context database.UserDatabaseAdapter) *AuthorizationControllers {
	controllers := AuthorizationControllers{
		UserAdapter: context,
	}
	return &controllers
}
