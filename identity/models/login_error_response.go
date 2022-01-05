package models

// LoginErrorResponse entity
type LoginErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
