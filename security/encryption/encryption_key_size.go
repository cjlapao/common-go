package encryption

import (
	"bytes"
	"encoding/json"
)

type EncryptionKeySize int64

const (
	Bit256  EncryptionKeySize = 256
	Bit384  EncryptionKeySize = 384
	Bit512  EncryptionKeySize = 512
	Bit1024 EncryptionKeySize = 1024
	Bit2048 EncryptionKeySize = 2048
	Bit4096 EncryptionKeySize = 4096
)

func (s EncryptionKeySize) String() string {
	return toEncryptionKeySizeString[s]
}

func (s EncryptionKeySize) FromString(key string) EncryptionKeySize {
	return toEncryptionKeySizeID[key]
}

var toEncryptionKeySizeString = map[EncryptionKeySize]string{
	Bit256:  "256",
	Bit384:  "384",
	Bit512:  "512",
	Bit1024: "1024",
	Bit2048: "2048",
	Bit4096: "4096",
}

var toEncryptionKeySizeID = map[string]EncryptionKeySize{
	"256":  Bit256,
	"384":  Bit384,
	"512":  Bit512,
	"1024": Bit1024,
	"2048": Bit2048,
	"4096": Bit4096,
}

func (s EncryptionKeySize) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toEncryptionKeySizeString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *EncryptionKeySize) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = toEncryptionKeySizeID[j]
	return nil
}
