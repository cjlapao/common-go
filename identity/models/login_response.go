package models

import "time"

// LoginResponse entity
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresIn    time.Time `json:"expires_in"`
	CreatedAt    time.Time `json:"created_at"`
}
