package jwt

type RawCertificateHeader struct {
	KeyId string `json:"kid,omitempty"`
	X5T   string `json:"x5t,omitempty"`
}
