package identity

// LoginResponse entity
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Expiring    string `json:"expiring"`
}
