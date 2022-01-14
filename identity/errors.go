package identity

import "github.com/cjlapao/common-go/identity/models"

// Error list for the identity package
// it uses the following framework identification
// packageId.errorCode
var (
	// ErrNoContext No User database context found error response
	ErrNoContext = models.NewOAuthErrorResponse(models.OAuthInvalidRequestError, "Context is null and cannot be validated")
)
