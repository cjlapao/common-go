package jwt_keyvault

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cjlapao/common-go/identity/jwk"
	"github.com/cjlapao/common-go/security/encryption"
)

type JwtKeyVaultItem struct {
	ID                string
	Type              encryption.EncryptionKey
	Size              encryption.EncryptionKeySize
	Thumbprint        string
	Certificate       *x509.Certificate
	EncodedPrivateKey string
	PrivateKey        interface{}
	EncodedPublicKey  string
	PublicKey         interface{}
	IsDefault         bool
	JWK               *jwk.JsonWebKeys
}

type JwtKeyVaultService struct {
	Keys []*JwtKeyVaultItem
}

var globalKeyVault *JwtKeyVaultService

func NewKeyVault() *JwtKeyVaultService {
	keyvault := JwtKeyVaultService{
		Keys: make([]*JwtKeyVaultItem, 0),
	}

	globalKeyVault = &keyvault

	return globalKeyVault
}

func Get() *JwtKeyVaultService {
	if globalKeyVault != nil {
		return globalKeyVault
	}

	return NewKeyVault()
}

func (kv *JwtKeyVaultService) WithCertificate(certificate x509.Certificate, privateKey interface{}) *JwtKeyVaultService {
	return kv
}

func (kv *JwtKeyVaultService) WithBase64RsaKey(id string, privateKey string) *JwtKeyVaultService {
	private := encryption.RSAHelper{}.DecodePrivateKeyFromBase64Pem(privateKey)
	return kv.WithRsaKey(id, private)
}

func (kv *JwtKeyVaultService) WithRsaKey(id string, privateKey *rsa.PrivateKey) *JwtKeyVaultService {
	if !kv.keyExists(id) {
		key := JwtKeyVaultItem{
			ID:         id,
			PrivateKey: privateKey,
			PublicKey:  privateKey.PublicKey,
		}
		key.Type = key.Type.FromString(fmt.Sprintf("RS%v", privateKey.Size()))
		key.Size = key.Size.FromString(fmt.Sprintf("%v", privateKey.Size()))
		key.Thumbprint = encryption.GetBase64KeyFingerprint(privateKey)

		x509PrivateEncodedBlock := x509.MarshalPKCS1PrivateKey(privateKey)
		x509PublicEncodedBlock := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
		key.EncodedPrivateKey = base64.StdEncoding.EncodeToString(x509PrivateEncodedBlock)
		key.EncodedPublicKey = base64.StdEncoding.EncodeToString(x509PublicEncodedBlock)

		key.JWK = jwk.New()
		key.JWK.Add(id, privateKey)

		if len(kv.Keys) == 0 {
			key.IsDefault = true
		}

		kv.Keys = append(kv.Keys, &key)
	}
	return kv
}

func (kv *JwtKeyVaultService) WithBase64EcdsaKey(id string, privateKey string) *JwtKeyVaultService {
	private := encryption.ECDSAHelper{}.DecodePrivateKeyFromBase64Pem(privateKey)
	return kv.WithEcdsaKey(id, private)
}

func (kv *JwtKeyVaultService) WithEcdsaKey(id string, privateKey *ecdsa.PrivateKey) *JwtKeyVaultService {
	if !kv.keyExists(id) {
		key := JwtKeyVaultItem{
			ID:         id,
			PrivateKey: privateKey,
			PublicKey:  privateKey.PublicKey,
		}
		key.Type = key.Type.FromString("ES" + fmt.Sprintf("%v", privateKey.Params().BitSize))
		key.Size = key.Size.FromString(fmt.Sprintf("%v", privateKey.Params().BitSize))
		key.Thumbprint = encryption.GetBase64KeyFingerprint(privateKey)

		x509PrivateEncodedBlock, _ := x509.MarshalECPrivateKey(privateKey)
		genericPublicKey, _ := x509.MarshalPKIXPublicKey(privateKey.PublicKey)
		key.EncodedPrivateKey = base64.StdEncoding.EncodeToString(x509PrivateEncodedBlock)
		key.EncodedPublicKey = base64.StdEncoding.EncodeToString(genericPublicKey)

		key.JWK = jwk.New()
		key.JWK.Add(id, privateKey)

		if len(kv.Keys) == 0 {
			key.IsDefault = true
		}

		kv.Keys = append(kv.Keys, &key)
	}
	return kv
}

func (kv *JwtKeyVaultService) WithBase64HmacKey(id string, privateKey string, size encryption.EncryptionKeySize) *JwtKeyVaultService {
	private, _ := base64.StdEncoding.DecodeString(privateKey)
	return kv.WithHmacKey(id, string(private), size)
}

func (kv *JwtKeyVaultService) WithHmacKey(id string, privateKey string, size encryption.EncryptionKeySize) *JwtKeyVaultService {
	if !kv.keyExists(id) {
		key := JwtKeyVaultItem{
			PrivateKey: privateKey,
		}
		key.Type = key.Type.FromString("HS" + size.String())
		key.Size = size
		key.Thumbprint = encryption.GetBase64KeyFingerprint(privateKey)

		key.EncodedPrivateKey = base64.StdEncoding.EncodeToString([]byte(privateKey))

		if len(kv.Keys) == 0 {
			key.IsDefault = true
		}

		kv.Keys = append(kv.Keys, &key)
	}
	return kv
}

func (kv *JwtKeyVaultService) SetDefaultKey(id string) {
	if kv.keyExists(id) {
		// Removing all defaults from other keys
		for _, key := range kv.Keys {
			key.IsDefault = false
		}

		for _, key := range kv.Keys {
			if strings.EqualFold(key.ID, id) {
				key.IsDefault = true
				break
			}
		}
	}
}

func (kv *JwtKeyVaultService) GetDefaultKey() *JwtKeyVaultItem {
	for _, key := range kv.Keys {
		if key.IsDefault {
			return key
		}
	}

	return nil
}

func (kv *JwtKeyVaultService) GetKey(id string) *JwtKeyVaultItem {
	for _, key := range kv.Keys {
		if strings.EqualFold(key.ID, id) {
			return key
		}
	}

	return nil
}

func (kv *JwtKeyVaultService) keyExists(id string) bool {
	for _, key := range kv.Keys {
		if strings.EqualFold(key.ID, id) {
			return true
		}
	}

	return false
}
