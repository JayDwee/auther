// Package auth contains functions for generating and verifying JWT tokens.
package jwt

import (
	"auther/internal/database/application"
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

// Constants for JWT token generation and verification.
const (
	// TokenExpiry is the expiry duration of the JWT token.
	TokenExpiry = 1 * time.Hour
)

//func GenerateJWK(as *application.Entity) (jwk.Set, error) {
//	key, err := jwk.New(jwk.NewRSAPrivateKey())
//	if err != nil {
//		return nil, fmt.Errorf("failed to generate JWK: %w", err)
//	}
//
//	as.JWKs.Add(key)
//	err = application.Repository.Save(context.TODO(), as)
//	if err != nil {
//		return nil, fmt.Errorf("failed to save JWK: %w", err)
//	}
//
//	return as.JWKs, nil
//}

// GenerateToken generates a JWT token with the given audience and scope.
// It returns the signed token and an error if any.
func GenerateToken(app *application.Entity, sub string, aud []string, iss string, scope string) (string, error) {
	token := jwt.New()

	for key, val := range map[string]interface{}{
		jwt.AudienceKey:   aud,
		jwt.SubjectKey:    sub,
		jwt.ExpirationKey: time.Now().Add(TokenExpiry).Unix(),
		jwt.IssuedAtKey:   time.Now().Unix(),
		jwt.IssuerKey:     iss,
		jwt.NotBeforeKey:  time.Now().Unix(),
		jwt.JwtIDKey:      uuid.New().String(),
		//jwt.SubjectKey: nil,
		"scope": scope} {
		err := token.Set(key, val)
		if err != nil {
			return "", fmt.Errorf("failed to set %s: %w", key, err)
		}
	}

	jwkSet, ok := app.JWKs.(jwk.Set)
	if !ok {
		return "", fmt.Errorf("JWKs is not a jwk.Set")
	}

	key, exists := jwkSet.LookupKeyID(app.ActiveKID)
	if !exists {
		return "", fmt.Errorf("JWK_KID not found in %s", key)
	}

	headers := jws.NewHeaders()
	err := headers.Set(jws.TypeKey, "at+JWT")
	if err != nil {
		return "", err
	}
	payload, err := jwt.NewSerializer().Serialize(token)
	signedKey, err := jws.Sign(payload, jwa.SignatureAlgorithm(key.Algorithm()), key, jws.WithHeaders(headers))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signedKey), nil
}
