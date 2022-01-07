package models

type OicdIntrospectResponse struct {
	ID        string `json:"jti,omitempty"`
	Active    bool   `json:"active"`
	TokenType string `json:"token_type,omitempty"`
	Subject   string `json:"sub,omitempty"`
	ClientId  string `json:"client_id,omitempty"`
	ExpiresAt string `json:"exp,omitempty"`
	IssuedAt  string `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
}
