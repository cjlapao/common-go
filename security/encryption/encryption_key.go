package encryption

import (
	"bytes"
	"encoding/json"
)

type EncryptionKey int64

const (
	HS256 EncryptionKey = iota
	HS384
	HS512
	RS256
	RS384
	RS512
	ES256
	ES384
	ES512
	PS256
	PS384
	PS512
)

func (j EncryptionKey) String() string {
	return toEncryptionKeyString[j]
}

func (j EncryptionKey) FromString(keyType string) EncryptionKey {
	return toEncryptionKeyID[keyType]
}

func (j EncryptionKey) GetFamily() EncryptionKeyType {
	return toEncryptionKeyFamily[j]
}

func (s EncryptionKey) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toEncryptionKeyString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *EncryptionKey) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = toEncryptionKeyID[j]
	return nil
}

var toEncryptionKeyString = map[EncryptionKey]string{
	HS256: "HS256",
	HS384: "HS384",
	HS512: "HS512",
	RS256: "RS256",
	RS384: "RS384",
	RS512: "RS512",
	ES256: "ES256",
	ES384: "ES384",
	ES512: "ES512",
	PS256: "PS256",
	PS384: "PS384",
	PS512: "PS512",
}

var toEncryptionKeyID = map[string]EncryptionKey{
	"HS256": HS256,
	"HS384": HS384,
	"HS512": HS512,
	"RS256": RS256,
	"RS384": RS384,
	"RS512": RS512,
	"ES256": ES256,
	"ES384": ES384,
	"ES512": ES512,
	"P256":  PS256,
	"PS384": PS384,
	"PS512": PS512,
}

var toEncryptionKeyFamily = map[EncryptionKey]EncryptionKeyType{
	HS256: HMAC,
	HS384: HMAC,
	HS512: HMAC,
	RS256: RSA,
	RS384: RSA,
	RS512: RSA,
	ES256: ECDSA,
	ES384: ECDSA,
	ES512: ECDSA,
	PS256: RSA,
	PS384: RSA,
	PS512: RSA,
}
