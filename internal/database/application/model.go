// Package application
/*
{
	"id": "kb2",
	"resource_owner_url": "https://login.koalabot.uk/",
	"clients": [
		{
			"client_id": "kb2",
			"client_secret": "secret", // Encrypted
			"audiences": [
				"https://kb2.koalabot.uk/api"
			]
		}
	],
	"custom_urls": [
		"https://auther.koalabot.uk"
	],
	"token_type": "JWE", // JWS, JWE, Opaque
	"active_kid": "78921401249908512",
	"jwks": "712bh41bh5gy12" // Encrypted
}
*/
package application

type Client struct {
	ClientId     string   `dynamodbav:"client_id"`
	ClientSecret string   `dynamodbav:"client_secret"` // TODO: Encrypt
	Audiences    []string `dynamodbav:"audiences"`
}

// Entity

type Entity struct {
	Id               string      `dynamodbav:"id"`
	ResourceOwnerUrl string      `dynamodbav:"resource_owner_url"`
	Clients          []Client    `dynamodbav:"clients"`
	CustomUrls       []string    `dynamodbav:"custom_urls"`
	ActiveKID        string      `dynamodbav:"active_kid"`
	JWKs             interface{} `dynamodbav:"jwks"`
}
