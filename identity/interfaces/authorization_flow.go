package interfaces

import "github.com/cjlapao/common-go/identity/models"

type AuthorizationFlow interface {
	Authorize(request *models.OAuthLoginRequest, tenant string) (*models.OAuthLoginResponse, *models.OAuthErrorResponse)
	RefreshToken(request *models.OAuthLoginRequest, tenant string) (*models.OAuthLoginResponse, *models.OAuthErrorResponse)
	ValidateEmailToken(request *models.OAuthLoginRequest, tenant string) (*models.OAuthLoginResponse, *models.OAuthErrorResponse)
}
