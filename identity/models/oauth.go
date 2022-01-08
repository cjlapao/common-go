package models

import (
	"bytes"
	"encoding/json"
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
