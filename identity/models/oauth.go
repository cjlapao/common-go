package models

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/security/encryption"
)

// OAuthGrantType Enum
type OAuthGrantType int64

const (
	OAuthPasswordGrant OAuthGrantType = iota
)

func (oauthGrantType OAuthGrantType) String() string {
	return toOAuthGrantTypeString[oauthGrantType]
}

func (oauthGrantType OAuthGrantType) FromString(keyType string) OAuthGrantType {
	return toOAuthGrantTypeID[keyType]
}

var toOAuthGrantTypeString = map[OAuthGrantType]string{
	OAuthPasswordGrant: "password",
}

var toOAuthGrantTypeID = map[string]OAuthGrantType{
	"password": OAuthPasswordGrant,
}

func (oauthGrantType OAuthGrantType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toOAuthGrantTypeString[oauthGrantType])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (oauthGrantType *OAuthGrantType) UnmarshalJSON(b []byte) error {
	var key string
	err := json.Unmarshal(b, &key)
	if err != nil {
		return err
	}

	*oauthGrantType = toOAuthGrantTypeID[key]
	return nil
}

// OAuthErrorType Enum
type OAuthErrorType int64

const (
	OAuthInvalidRequestError OAuthErrorType = iota
	OAuthInvalidClientError
	OAuthInvalidGrant
	OAuthInvalidScope
	OAuthUnauthorizedClient
	OAuthUnsupportedGrantType
)

func (oAuthErrorType OAuthErrorType) String() string {
	return toOAuthErrorTypeString[oAuthErrorType]
}

func (oAuthErrorType OAuthErrorType) FromString(keyType string) OAuthErrorType {
	return toOAuthErrorTypeID[keyType]
}

var toOAuthErrorTypeString = map[OAuthErrorType]string{
	OAuthInvalidRequestError:  "invalid_request",
	OAuthInvalidClientError:   "invalid_client",
	OAuthInvalidGrant:         "invalid_grant",
	OAuthInvalidScope:         "invalid_scope",
	OAuthUnauthorizedClient:   "unauthorized_client",
	OAuthUnsupportedGrantType: "unsupported_grant_type",
}

var toOAuthErrorTypeID = map[string]OAuthErrorType{
	"invalid_request":        OAuthInvalidRequestError,
	"invalid_client":         OAuthInvalidClientError,
	"invalid_grant":          OAuthInvalidGrant,
	"invalid_scope":          OAuthInvalidScope,
	"unauthorized_client":    OAuthUnauthorizedClient,
	"unsupported_grant_type": OAuthUnsupportedGrantType,
}

func (oAuthErrorType OAuthErrorType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toOAuthErrorTypeString[oAuthErrorType])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (oAuthErrorType *OAuthErrorType) UnmarshalJSON(b []byte) error {
	var key string
	err := json.Unmarshal(b, &key)
	if err != nil {
		return err
	}

	*oAuthErrorType = toOAuthErrorTypeID[key]
	return nil
}

// OAuthLoginRequest Entity
type OAuthLoginRequest struct {
	GrantType    string `json:"grant_type"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

// OAuthLoginRequest Entity
type OAuthLoginResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// OAuthIntrospectResponse entity
type OAuthIntrospectResponse struct {
	ID        string `json:"jti,omitempty"`
	Active    bool   `json:"active"`
	TokenType string `json:"token_type,omitempty"`
	Subject   string `json:"sub,omitempty"`
	ClientId  string `json:"client_id,omitempty"`
	ExpiresAt string `json:"exp,omitempty"`
	IssuedAt  string `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
}

// OAuthErrorResponse entity
type OAuthErrorResponse struct {
	Error            OAuthErrorType `json:"error"`
	ErrorDescription string         `json:"error_description,omitempty"`
	ErrorUri         string         `json:"error_uri,omitempty"`
}

func NewOAuthErrorResponse(err OAuthErrorType, description string) OAuthErrorResponse {
	errorResponse := OAuthErrorResponse{
		Error:            err,
		ErrorDescription: description,
	}
	return errorResponse

}

func (err OAuthErrorResponse) String() string {
	return fmt.Sprintf("An error occurred, %v: %v", err.Error.String(), err.ErrorDescription)
}

func (err OAuthErrorResponse) Log(extraLogs ...string) {
	logger := log.Get()
	if len(extraLogs) == 0 {
		logger.Error(err.String())
	} else {
		logger.Error("%v %v", err.String(), extraLogs[0])
	}
}

type OAuthConfigurationResponse struct {
	Issuer                             string   `json:"issuer"`
	JwksURI                            string   `json:"jwks_uri"`
	AuthorizationEndpoint              string   `json:"authorization_endpoint"`
	TokenEndpoint                      string   `json:"token_endpoint"`
	UserinfoEndpoint                   string   `json:"userinfo_endpoint"`
	EndSessionEndpoint                 string   `json:"end_session_endpoint"`
	CheckSessionIframe                 string   `json:"check_session_iframe"`
	RevocationEndpoint                 string   `json:"revocation_endpoint"`
	IntrospectionEndpoint              string   `json:"introspection_endpoint"`
	DeviceAuthorizationEndpoint        string   `json:"device_authorization_endpoint"`
	FrontchannelLogoutSupported        bool     `json:"frontchannel_logout_supported"`
	FrontchannelLogoutSessionSupported bool     `json:"frontchannel_logout_session_supported"`
	BackchannelLogoutSupported         bool     `json:"backchannel_logout_supported"`
	BackchannelLogoutSessionSupported  bool     `json:"backchannel_logout_session_supported"`
	ScopesSupported                    []string `json:"scopes_supported"`
	ClaimsSupported                    []string `json:"claims_supported"`
	GrantTypesSupported                []string `json:"grant_types_supported"`
	ResponseTypesSupported             []string `json:"response_types_supported"`
	ResponseModesSupported             []string `json:"response_modes_supported"`
	TokenEndpointAuthMethodsSupported  []string `json:"token_endpoint_auth_methods_supported"`
	SubjectTypesSupported              []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported   []string `json:"id_token_signing_alg_values_supported"`
	CodeChallengeMethodsSupported      []string `json:"code_challenge_methods_supported"`
	RequestParameterSupported          bool     `json:"request_parameter_supported"`
}

type OAuthJwksResponse struct {
	Keys []OAuthJwksKey `json:"keys"`
}

type OAuthJwksKey struct {
	ID              string                       `json:"kid"`
	Algorithm       encryption.EncryptionKey     `json:"alg"`
	AlgorithmFamily encryption.EncryptionKeyType `json:"kty"`
	Use             string                       `json:"use"`
	X5C             []string                     `json:"x5c"`
	Exponent        string                       `json:"e,omitempty"`
	Modulus         string                       `json:"n,omitempty"`
	Curve           string                       `json:"curve,omitempty"`
	X               string                       `json:"x,omitempty"`
	Y               string                       `json:"y,omitempty"`
	Thumbprint      string                       `json:"x5t"`
}

type OAuthRegisterRequest struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Email     string   `json:"email"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Roles     []string `json:"roles"`
	Claims    []string `json:"claims"`
}

type OAuthRevokeRequest struct {
	ClientID  string `json:"client_id"`
	GrantType string `json:"grant_type"`
}
