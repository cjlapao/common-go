package jwt

type RawCertificateHeader struct {
	Algorithm string `json:"alg,omitempty"`
	KeyId     string `json:"kid,omitempty"`
	X5T       string `json:"x5t,omitempty"`
}
