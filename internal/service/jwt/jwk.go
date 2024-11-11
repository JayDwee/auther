package jwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
)

func GenerateJWK(sig jwa.SignatureAlgorithm) (key jwk.Key, err error) {
	var privateKey interface{}
	switch sig {
	case jwa.ES256:
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case jwa.ES256K:
		privateKey, err = ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	case jwa.ES384:
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case jwa.ES512:
		privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	case jwa.EdDSA:
		_, privateKey, err = ed25519.GenerateKey(rand.Reader)
	case jwa.HS256:
		tmpKey := make([]byte, 32)
		_, err = rand.Read(tmpKey)
		privateKey = tmpKey
	case jwa.HS384:
		tmpKey := make([]byte, 48)
		_, err = rand.Read(tmpKey)
		privateKey = tmpKey
	case jwa.HS512:
		tmpKey := make([]byte, 64)
		_, err = rand.Read(tmpKey)
		privateKey = tmpKey
	case jwa.NoSignature:
		return nil, nil
	case jwa.PS256, jwa.RS256:
		privateKey, err = rsa.GenerateKey(rand.Reader, 256*8)
	case jwa.PS384, jwa.RS384:
		privateKey, err = rsa.GenerateKey(rand.Reader, 384*8)
	case jwa.PS512, jwa.RS512:
		privateKey, err = rsa.GenerateKey(rand.Reader, 512*8)
	default:
		return nil, fmt.Errorf(`invalid signature algorithm %s`, sig)
	}
	if err != nil {
		return nil, fmt.Errorf(`failed to generate key`, err)
	}

	key, err = jwk.New(privateKey)
	if err != nil {
		return
	}

	err = jwk.AssignKeyID(key)
	if err != nil {
		return
	}

	err = key.Set(jwk.AlgorithmKey, sig)
	if err != nil {
		return
	}

	return
}
