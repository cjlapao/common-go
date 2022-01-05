package models

import "time"

type UserToken struct {
	Token        string
	NotBefore    time.Time
	ExpiresAt    time.Time
	Audiences    []string
	Issuer       string
	UsedKeyID    string
	RefreshToken string
}
