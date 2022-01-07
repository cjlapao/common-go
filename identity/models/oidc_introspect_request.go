package models

type OicdIntrospectRequest struct {
	Token         string `json:"token"`
	TokenTypeHint string `json:"token_type_hint"`
	ClientId      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
}
