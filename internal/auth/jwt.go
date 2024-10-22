// Package auth contains functions for generating and verifying JWT tokens.
package auth

import (
	"auther/internal/auth_server"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

// Constants for JWT token generation and verification.
const (
	// TokenExpiry is the expiry duration of the JWT token.
	TokenExpiry = 1 * time.Hour

	// S3Key is the path to the private JWK file.
	//S3Key = "private/jwks.json"
)

var (
	// Issuer is the issuer of the JWT token.
	Issuer = os.Getenv("ISSUER")

	//// Default Client TODO Remove, custom per client
	//DefaultClient = os.Getenv("DEFAULT_CLIENT")
	//
	//// S3Bucket S3 bucket containing the private JWK file.
	//S3Bucket = os.Getenv("S3_BUCKET") // TODO Write to S3 when new keys generated
)

//
//// jwkSet is the set of JWKs used for signing and verifying tokens.
//var jwkSet jwk.Set
//
//// init loads the private JWK set from S3.
//func init() {
//	var err error
//	jwkSet, err = LoadPrivateJWKFromDynamo()
//	if err != nil {
//		log.Fatalf("failed to load private JWK: %v", err)
//	}
//}
//
//func LoadPrivateJWKFromDynamo() (jwk.Set, error) {
//	key, err := auth_server.AuthorizationServerRepository.GetByHashKey(context.TODO(), DefaultClient)
//	if err != nil {
//		return nil, err
//	}
//
//	return key.JWKs, nil
//}
//
//// LoadPrivateJWK loads the private JWK set from S3.
//// It returns the set of JWKs and an error if any.
//func LoadPrivateJWKFromS3() (jwk.Set, error) {
//	log.Println("Starting to load private JWK from S3")
//	cfg, err := config.LoadDefaultConfig(context.TODO())
//	if err != nil {
//		return nil, fmt.Errorf("failed to load AWS config: %w", err)
//	}
//
//	client := s3.NewFromConfig(cfg)
//	log.Println("Connected to client")
//
//	output, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
//		Bucket: aws.String(S3Bucket),
//		Key:    aws.String(S3Key),
//	})
//	if err != nil {
//		return nil, fmt.Errorf("failed to get S3 object: %w", err)
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Fatalf("failed to close S3 object body: %v", err)
//		}
//	}(output.Body)
//
//	body, err := io.ReadAll(output.Body)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read S3 object body: %w", err)
//	}
//	log.Println("S3 read complete")
//
//	jwkSet, err := jwk.Parse(body)
//	if err != nil {
//		return nil, fmt.Errorf("failed to parse JWK: %w", err)
//	}
//
//	return jwkSet, nil
//}

// GenerateToken generates a JWT token with the given audience and scope.
// It returns the signed token and an error if any.
func GenerateToken(as *auth_server.AuthorizationServerModel, audience string, scope string) (string, error) {
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

	jwkRsa, exists := as.JWKs.LookupKeyID(as.ActiveKID)
	if !exists {
		return "", fmt.Errorf("JWK_KID not found in %s", jwkRsa)
	}

	signedKey, err := jwt.NewSerializer().Sign(jwa.RS512, jwkRsa).Serialize(token)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signedKey), nil
}
