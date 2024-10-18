// Package auth contains functions for generating and verifying JWT tokens.
package auth

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// Constants for JWT token generation and verification.
const (
	// Issuer is the issuer of the JWT token.
	Issuer = "temp.auther.koalabot.uk"
	// TokenExpiry is the expiry duration of the JWT token.
	TokenExpiry = 1 * time.Hour

	// S3Bucket S3 bucket containing the private JWK file.
	S3Bucket = "auther"
	// S3Key is the path to the private JWK file.
	S3Key = "private/jwks.json"
)

// jwkSet is the set of JWKs used for signing and verifying tokens.
var jwkSet jwk.Set

// init loads the private JWK set from S3.
func init() {
	var err error
	jwkSet, err = loadPrivateJWK()
	if err != nil {
		log.Fatalf("failed to load private JWK: %v", err)
	}
}

// loadPrivateJWK loads the private JWK set from S3.
// It returns the set of JWKs and an error if any.
func loadPrivateJWK() (jwk.Set, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	output, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(S3Key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get S3 object: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("failed to close S3 object body: %v", err)
		}
	}(output.Body)

	body, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read S3 object body: %w", err)
	}

	jwkSet, err := jwk.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWK: %w", err)
	}

	return jwkSet, nil
}

// GenerateToken generates a JWT token with the given audience and scope.
// It returns the signed token and an error if any.
func GenerateToken(audience, scope string) (string, error) {
	token := jwt.New()

	for key, val := range map[string]interface{}{
		jwt.AudienceKey:   audience,
		jwt.ExpirationKey: time.Now().Add(TokenExpiry).Unix(),
		jwt.IssuedAtKey:   time.Now().Unix(),
		jwt.IssuerKey:     Issuer,
		jwt.NotBeforeKey:  time.Now().Unix(),
		jwt.JwtIDKey:      uuid.New().String(),
		//jwt.SubjectKey: nil,
		"scope": scope} {
		err := token.Set(key, val)
		if err != nil {
			return "", fmt.Errorf("failed to set %s: %w", key, err)
		}
	}

	jwkRsa, exists := jwkSet.LookupKeyID(os.Getenv("JWK_KID"))
	if !exists {
		return "", fmt.Errorf("JWK_KID not found in %s", jwkRsa)
	}

	signedKey, err := jwt.NewSerializer().Sign(jwa.RS512, jwkRsa).Serialize(token)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signedKey), nil
}
