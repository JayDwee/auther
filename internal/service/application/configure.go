package application

import (
	"auther/internal/database/application"
	"auther/internal/service/jwt"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"os"
)

func UpdateJWKS3(app *application.Entity) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	keySet, ok := app.JWKs.(jwk.Set)
	if !ok {
		return fmt.Errorf("failed to cast JWKs to jwk.Set")
	}
	keySet, err = keySet.Clone()
	keySetIterator := keySet.Iterate(context.TODO())

	publicKeys := jwk.NewSet()

	for keySetIterator.Next(context.TODO()) {
		key := keySetIterator.Pair().Value.(jwk.Key)
		publicKey, err := key.PublicKey()
		if err != nil {
			return err
		}
		publicKeys.Add(publicKey)
	}

	client := s3.NewFromConfig(cfg)
	marshal, err := json.Marshal(publicKeys)
	if err != nil {
		return err
	}

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("auther-" + os.Getenv("DEPLOYMENT_ENV")),
		Key:    aws.String("applications/" + app.Id + "/.well-known/jwks.json"),
		Body:   bytes.NewReader(marshal),
	})
	if err != nil {
		return err
	}
	return nil
}
func Create(applicationId string) (entity *application.Entity, err error) {
	// Create Applicaiton in DB
	entity = &application.Entity{
		Id:               applicationId,
		ResourceOwnerUrl: "",
		Clients:          nil,
		CustomUrls:       nil,
		ActiveKID:        "",
		JWKs:             nil,
	}

	err = application.Repository.Save(context.TODO(), entity)
	if err != nil {
		return
	}

	return
}

func AddGeneratedKey(app *application.Entity, alg jwa.SignatureAlgorithm) (err error) {
	newKey, err := jwt.GenerateJWK(alg)
	if err != nil {
		return
	}
	if app.JWKs == nil {
		app.JWKs = jwk.NewSet()
	}

	jwkSet, ok := app.JWKs.(jwk.Set)
	if !ok {
		return fmt.Errorf("failed to cast JWKs to jwk.Set")
	}
	jwkSet.Add(newKey)
	app.ActiveKID = newKey.KeyID()
	err = application.Repository.Save(context.TODO(), app)
	if err != nil {
		return
	}

	err = UpdateJWKS3(app)
	if err != nil {
		return
	}
	return
}

func RemoveKey(app *application.Entity, kid string) (err error) {
	jwkSet, ok := app.JWKs.(jwk.Set)
	if !ok {
		return fmt.Errorf("failed to cast JWKs to jwk.Set")
	}
	key, found := jwkSet.LookupKeyID(kid)
	if !found {
		return fmt.Errorf("key not found")
	}

	deleted := jwkSet.Remove(key)
	if !deleted {
		return fmt.Errorf("key not deleted")
	}

	err = application.Repository.Save(context.TODO(), app)
	if err != nil {
		return err
	}

	return
}
