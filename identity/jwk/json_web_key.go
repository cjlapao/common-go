package jwk

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/cjlapao/common-go/guard"
	"github.com/cjlapao/common-go/security/encryption"
)

type JsonWebKey struct {
	ID              string                       `json:"kid"`
	Algorithm       encryption.EncryptionKey     `json:"alg"`
	AlgorithmFamily encryption.EncryptionKeyType `json:"kty"`
	Use             string                       `json:"use"`
	X5C             []string                     `json:"x5c"`
	Exponent        string                       `json:"e,omitempty"`
	Modulus         string                       `json:"n,omitempty"`
	Curve           string                       `json:"curve,omitempty"`
	X               string                       `json:"x,omitempty"`
	Y               string                       `json:"y,omitempty"`
	Thumbprint      string                       `json:"x5t"`
}

func NewKey(privateKey interface{}) *JsonWebKey {
	return NewKeyWithId("", privateKey)
}

func NewKeyWithId(id string, privateKey interface{}) *JsonWebKey {
	guard.FatalEmptyOrNil(privateKey)

	key := JsonWebKey{
		Use: "sig",
		X5C: make([]string, 0),
	}

	fingerprint := encryption.GetSHA1KeyFingerprint(privateKey)
	key.Thumbprint = base64.URLEncoding.EncodeToString(fingerprint[:])
	if id == "" {
		key.ID = key.Thumbprint
	}

	switch kt := privateKey.(type) {
	case *rsa.PrivateKey:
		publicKey := kt.PublicKey
		key.Algorithm = key.Algorithm.FromString("RS" + fmt.Sprintf("%v", kt.Size()))
		key.AlgorithmFamily = key.Algorithm.GetFamily()
		exponentBytes := new(big.Int).SetInt64(int64(publicKey.E))
		key.Exponent = base64.URLEncoding.EncodeToString(exponentBytes.Bytes())
		key.Modulus = base64.URLEncoding.EncodeToString(publicKey.N.Bytes())
		encodedPublicKey := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&publicKey))
		key.X5C = append(key.X5C, encodedPublicKey)
	case *ecdsa.PrivateKey:
		publicKey := kt.PublicKey
		key.Algorithm = key.Algorithm.FromString("ES" + fmt.Sprintf("%v", publicKey.Params().BitSize))
		key.AlgorithmFamily = key.Algorithm.GetFamily()
		key.Curve = publicKey.Params().Name
		key.X = base64.URLEncoding.EncodeToString(publicKey.X.Bytes())
		key.Y = base64.URLEncoding.EncodeToString(publicKey.Y.Bytes())
		genericKey, _ := x509.MarshalPKIXPublicKey(&publicKey)
		encodedPublicKey := base64.StdEncoding.EncodeToString(genericKey)
		key.X5C = append(key.X5C, encodedPublicKey)
	}

	return &key
}

func (k JsonWebKey) Validate(key interface{}) bool {
	switch kt := key.(type) {
	case rsa.PublicKey:
		return k.validateRsa(kt)
	case *rsa.PublicKey:
		return k.validateRsa(*kt)
	case rsa.PrivateKey:
		return k.validateRsa(kt.PublicKey)
	case *rsa.PrivateKey:
		return k.validateRsa(kt.PublicKey)
	case ecdsa.PublicKey:
		return k.validateEcdsa(kt)
	case *ecdsa.PublicKey:
		return k.validateEcdsa(*kt)
	case ecdsa.PrivateKey:
		return k.validateEcdsa(kt.PublicKey)
	case *ecdsa.PrivateKey:
		return k.validateEcdsa(kt.PublicKey)
	default:
		return false
	}
}

func Decode(jwk string) (*JsonWebKey, error) {
	var decodedKey JsonWebKey
	err := json.Unmarshal([]byte(jwk), &decodedKey)
	if err != nil {
		return nil, err
	}

	return &decodedKey, nil
}

func (k JsonWebKey) GetKey() interface{} {
	return k.GetPublicKeyAt(0)
}

func (k JsonWebKey) GetPublicKeyAt(index int) interface{} {
	if len(k.X5C) == 0 || len(k.X5C) < index {
		return nil
	}

	stringkey := k.X5C[index]
	key, err := base64.StdEncoding.DecodeString(stringkey)
	if err != nil {
		return nil
	}

	var publicKey interface{}
	switch k.AlgorithmFamily {
	case encryption.RSA:
		publicKey, err = x509.ParsePKCS1PublicKey(key)
		if err != nil {
			return nil
		}
		if !k.verifyKeyParameters(publicKey) {
			return nil
		}
		return publicKey
	case encryption.ECDSA:
		generickPublicKey, err := x509.ParsePKIXPublicKey(key)
		publicKey = generickPublicKey.(*ecdsa.PublicKey)
		if err != nil {
			return nil
		}
		if !k.verifyKeyParameters(publicKey) {
			return nil
		}
		return publicKey
	}

	return nil
}

func (k JsonWebKey) verifyKeyParameters(key interface{}) bool {
	switch kt := key.(type) {
	case rsa.PublicKey:
		return k.validateRsa(kt)
	case *rsa.PublicKey:
		return k.validateRsa(*kt)
	case ecdsa.PublicKey:
		return k.validateEcdsa(kt)
	case *ecdsa.PublicKey:
		return k.validateEcdsa(*kt)
	default:
		return false
	}
}

func (k JsonWebKey) validateRsa(publicKey rsa.PublicKey) bool {
	decodedExponent, _ := base64.URLEncoding.DecodeString(k.Exponent)
	decodedModulus, _ := base64.URLEncoding.DecodeString(k.Modulus)
	exponent := new(big.Int).SetBytes(decodedExponent)
	modulus := new(big.Int).SetBytes(decodedModulus)
	if publicKey.E == int(exponent.Int64()) && publicKey.N.Cmp(modulus) == 0 {
		return true
	}

	return false
}

func (k JsonWebKey) validateEcdsa(publicKey ecdsa.PublicKey) bool {
	decodedX, _ := base64.URLEncoding.DecodeString(k.X)
	decodedY, _ := base64.URLEncoding.DecodeString(k.Y)
	x := new(big.Int).SetBytes(decodedX)
	y := new(big.Int).SetBytes(decodedY)

	if publicKey.X.Cmp(x) == 0 && publicKey.Y.Cmp(y) == 0 {
		return true
	}

	return false
}
