package models

// LoginErrorResponse entity
type LoginErrorResponse struct {
	Code    string `json:"code" yaml:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
