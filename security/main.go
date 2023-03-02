package security

import (
	"crypto/sha256"
	"encoding/hex"

	log "github.com/cjlapao/common-go-logger"
)

var logger = log.Get()

// SHA256Encode Hash string with SHA256
func SHA256Encode(value string) string {
	hasher := sha256.New()
	bytes := []byte(value)
	hasher.Write(bytes)

	return hex.EncodeToString(hasher.Sum(nil))
}
