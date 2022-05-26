package templateutils

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rc4"
	"crypto/rsa"
	"crypto/sha256"
)

type (
	// rc4-
	rc4NS struct{}
	// des-
	desNS struct{}
	// des3-
	tripleDESNS struct{}

	// aes-{128,192,256}-{gcm,cbc,cfb,ofb}
	aesNS struct{}

	rsaNS       struct{}
	secp224r1NS struct{}
	eddsaNS     struct{}
)

func (rc4NS) Decrypt(key Bytes, args ...any) (ret string, err error) {
	data, err := toBytes(key)
	if err != nil {
		return
	}
	rc4.NewCipher(data)
	return
}

func (desNS) Decrypt(key Bytes, args ...any) (ret string, err error) {
	data, err := toBytes(key)
	if err != nil {
		return
	}
	des.NewCipher(data)
	return
}

func (tripleDESNS) Decrypt(key Bytes, args ...any) (ret string, err error) {
	data, err := toBytes(key)
	if err != nil {
		return
	}
	des.NewTripleDESCipher(data)
	return
}

func (aesNS) Decrypt(key Bytes, args ...any) (ret string, err error) {
	data, err := toBytes(key)
	if err != nil {
		return
	}
	b, err := aes.NewCipher(data)
	if err != nil {
		return
	}

	a, err := cipher.NewGCM(b)
	if err != nil {
		return
	}
	_ = a

	cipher.NewCFBEncrypter(b, nil)

	cipher.NewOFB(b, nil)
	return
}

func (rsaNS) Decrypt(k PrivateKey, args ...any) (ret string, err error) {
	rsa.DecryptOAEP(sha256.New(), rand.Reader, nil /* private key */, nil /* ciphertext */, nil /* label */)
	return
}

func (secp224r1NS) Decrypt(k PrivateKey, args ...any) (ret string, err error) {
	// switch t := toPrivateKey(k).(type) {
	// case ecdsa.PrivateKey:
	// case ed25519.PrivateKey:
	// }
	elliptic.GenerateKey(elliptic.P224(), rand.Reader)
	return
}

type PrivateKey any

func toPrivateKey(k PrivateKey) crypto.PrivateKey {
	// data := toBytes(k)
	// _ = data

	return nil
}
