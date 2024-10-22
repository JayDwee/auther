package auth_server

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/lestrrat-go/jwx/jwk"
	"log"
	"os"
)

var client *dynamodb.Client

func getClient() *dynamodb.Client {
	if client == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			panic(err)
		}

		client = dynamodb.NewFromConfig(cfg)
	}
	return client
}

type Client struct {
	ClientId     string `dynamodbav:"client_id"`
	ClientSecret string `dynamodbav:"client_secret"` // TODO: Encrypt
}

type AuthorizationServerModel struct {
	Id               string   `dynamodbav:"id"`
	ResourceOwnerUrl string   `dynamodbav:"resource_owner_url"`
	Clients          []Client `dynamodbav:"clients"`
	CustomUrls       []string `dynamodbav:"custom_urls"`
	ActiveKID        string   `dynamodbav:"active_kid"`
	JWKs             jwk.Set  `dynamodbav:"jwks"`
}

type authorizationServerEntity struct {
	*AuthorizationServerModel
	JWKs string `dynamodbav:"jwks"` // TODO: Encrypt
}

type authorizationServerRepository struct {
	tableName string
}

var AuthorizationServerRepository = authorizationServerRepository{
	tableName: "auther_" + os.Getenv("DEPLOYMENT_ENV") + "_authorization_server",
}

func (repo *authorizationServerRepository) Save(ctx context.Context, as *AuthorizationServerModel) error {
	asDto := authorizationServerEntity{
		AuthorizationServerModel: as,
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

	_, err = getClient().PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

func (repo *authorizationServerRepository) GetByHashKey(ctx context.Context, hashKey string) (as *AuthorizationServerModel, err error) {
	as = &AuthorizationServerModel{
		Id: hashKey,
	}
	asDto := authorizationServerEntity{
		AuthorizationServerModel: as,
	}
	idAv, err := attributevalue.Marshal(asDto.Id)

	response, err := getClient().GetItem(ctx, &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": idAv},
		TableName: aws.String(repo.tableName),
	})

	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", hashKey, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &asDto)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}
	// convert string to jwk.Set
	if asDto.JWKs != "" {
		jwks, err := jwk.Parse([]byte(asDto.JWKs))
		if err != nil {
			log.Printf("Couldn't parse jwks. Here's why: %v\n", err)
		}
		as.JWKs = jwks
	}

	return as, err
}
