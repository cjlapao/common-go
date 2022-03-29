package encryption

import (
	"bytes"
	"encoding/json"
)

type EncryptionKeyType int64

const (
	ECDSA EncryptionKeyType = iota
	HMAC
	RSA
)

func (j EncryptionKeyType) String() string {
	return toEncryptionKeyTypeString[j]
}

func (j EncryptionKeyType) FromString(keyType string) EncryptionKeyType {
	return toEncryptionKeyTypeID[keyType]
}

var toEncryptionKeyTypeString = map[EncryptionKeyType]string{
	RSA:   "RSA",
	ECDSA: "ECDSA",
	HMAC:  "HMAC",
}

var toEncryptionKeyTypeID = map[string]EncryptionKeyType{
	"RSA":   RSA,
	"ECDSA": ECDSA,
	"HMAC":  HMAC,
}

func (s EncryptionKeyType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toEncryptionKeyTypeString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *EncryptionKeyType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = toEncryptionKeyTypeID[j]
	return nil
}
