package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type UserToken struct {
	ID            string    `json:"jti,omitempty"`
	Scope         string    `json:"scope,omitempty"`
	User          string    `json:"sub,omitempty"`
	Issuer        string    `json:"iss,omitempty"`
	FirstName     string    `json:"given_name,omitempty"`
	LastName      string    `json:"family_name,omitempty"`
	UserID        string    `json:"uid,omitempty"`
	TenantId      string    `json:"tid,omitempty"`
	DisplayName   string    `json:"name,omitempty"`
	Email         string    `json:"email,omitempty"`
	EmailVerified bool      `json:"email_verified,omitempty"`
	Nonce         string    `json:"nonce,omitempty"`
	NotBefore     time.Time `json:"nbf,omitempty"`
	ExpiresAt     time.Time `json:"exp,omitempty"`
	IssuedAt      time.Time `json:"iat,omitempty"`
	Audiences     []string  `json:"aud,omitempty"`
	Roles         []string  `json:"roles,omitempty"`
	UsedKeyID     string    `json:"-"`
	RefreshToken  string    `json:"-"`
	Token         string    `json:"-"`
}

func (userToken *UserToken) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}
	m := f.(map[string]interface{})
	for k, v := range m {
		switch k {
		case "jti":
			userToken.ID = v.(string)
		case "scope":
			userToken.Scope = v.(string)
		case "sub":
			userToken.User = v.(string)
		case "iss":
			userToken.Issuer = v.(string)
		case "given_name":
			userToken.FirstName = v.(string)
		case "family_name":
			userToken.LastName = v.(string)
		case "email":
			userToken.Email = v.(string)
		case "email_verified":
			userToken.EmailVerified = v.(bool)
		case "uid":
			userToken.UserID = v.(string)
		case "tid":
			userToken.TenantId = v.(string)
		case "name":
			userToken.DisplayName = v.(string)
		case "nonce":
			userToken.Nonce = v.(string)
		case "nbf":
			userToken.NotBefore = time.Unix(int64(v.(float64)), 0)
		case "exp":
			userToken.ExpiresAt = time.Unix(int64(v.(float64)), 0)
		case "iat":
			userToken.IssuedAt = time.Unix(int64(v.(float64)), 0)
		case "aud":
			audienceValues := v.([]interface{})
			for _, v := range audienceValues {
				userToken.Audiences = append(userToken.Audiences, v.(string))
			}
		case "roles":
			rolesValues := v.([]interface{})
			for _, v := range rolesValues {
				userToken.Roles = append(userToken.Audiences, v.(string))
			}
		}
	}
	if userToken.DisplayName == "" {
		if userToken.FirstName != "" {
			userToken.DisplayName = userToken.FirstName
		}
		if userToken.LastName != "" {
			if userToken.DisplayName != "" {
				userToken.DisplayName = fmt.Sprintf("%v %v", userToken.DisplayName, userToken.LastName)
			} else {
				userToken.DisplayName = userToken.LastName
			}
		}
	}

	return nil
}
