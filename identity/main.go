package identity

import "github.com/cjlapao/common-go/log"

var logger = log.Get()

//TODO: Create API_KEY authorization
//TODO: Get jwt public key using openid configuration
//TODO: Cache the openid configuration for tokens based in the subject
//TODO: repurpose the JWT Keyvault as a generic key vault to keep secrets
//TODO: Make all errors variables for reusability purpose
//TODO: Move all log.error to log.exception for a cleaner implementation
