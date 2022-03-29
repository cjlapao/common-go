package jwt

import (
	"encoding/json"
	"fmt"

	"github.com/cjlapao/common-go/log"
	"github.com/pascaldekloe/jwt"
)

var logger = log.Get()

func GetTokenClaim(token string, claim string) string {
	if token == "" || claim == "" {
		return ""
	}

	jwtToken, err := jwt.ParseWithoutCheck([]byte(token))

	if err != nil {
		return ""
	}

	// Transforming token into a user token
	rawJsonToken, _ := jwtToken.Raw.MarshalJSON()
	var tokenMap map[string]interface{}
	err = json.Unmarshal(rawJsonToken, &tokenMap)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%v", tokenMap[claim])
}
