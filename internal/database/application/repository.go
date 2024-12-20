package application

import (
	"auther/internal/database"
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/lestrrat-go/jwx/jwk"
	"log"
	"os"
)

type internalApplicationEntity struct {
	*Entity
	JWKs string `dynamodbav:"jwks"` // TODO: Encrypt
}

type internalApplicationRepository struct {
	tableName string
}

func (repo *internalApplicationRepository) Save(ctx context.Context, as *Entity) error {
	asDto := internalApplicationEntity{
		Entity: as,
	}
	// convert jwk.Set to string
	if as.JWKs != nil {
		jwkString, err := json.Marshal(as.JWKs)
		if err != nil {
			return err
		}
		asDto.JWKs = string(jwkString)
	}

	item, err := attributevalue.MarshalMap(asDto)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	_, err = database.GetClient().PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

func (repo *internalApplicationRepository) GetByHashKey(ctx context.Context, hashKey string) (entity *Entity, err error) {
	entity = &Entity{
		Id: hashKey,
	}
	appDTO := internalApplicationEntity{
		Entity: entity,
	}
	idAv, err := attributevalue.Marshal(appDTO.Id)

	response, err := database.GetClient().GetItem(ctx, &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": idAv},
		TableName: aws.String(repo.tableName),
	})

	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", hashKey, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &appDTO)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}
	// convert string to jwk.Set
	if appDTO.JWKs != "" {
		jwks, err := jwk.Parse([]byte(appDTO.JWKs))
		if err != nil {
			log.Printf("Couldn't parse jwks. Here's why: %v\n", err)
		}
		entity.JWKs = jwks
	}

	return entity, err
}

func (repo *internalApplicationRepository) DeleteByHashKey(ctx context.Context, hashKey string) (err error) {

	_, err = database.GetClient().DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: hashKey},
		},
		TableName: aws.String(repo.tableName),
	})
	return err
}

var Repository = internalApplicationRepository{
	tableName: "auther_" + os.Getenv("DEPLOYMENT_ENV") + "_applications",
}
